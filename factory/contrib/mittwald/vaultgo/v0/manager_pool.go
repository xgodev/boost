package vault

import (
	"context"
	"fmt"

	"github.com/hashicorp/vault/api"
	vault "github.com/mittwald/vaultgo"
	"github.com/xgodev/boost/wrapper/log"
)

// ManagerPool configures managers from the secrets of one vault client.
type ManagerPool struct {
	managers []Manager
	client   *vault.Client
}

// NewManagerPool returns a new manager pool.
func NewManagerPool(client *vault.Client, managers ...Manager) *ManagerPool {
	return &ManagerPool{managers: managers, client: client}
}

// ManageAll builds a vault client with default options and configures the
// managers with it. Errors are logged.
func ManageAll(ctx context.Context, managers ...Manager) {
	client, err := NewClient(ctx)
	if err != nil {
		log.FromContext(ctx).Error(err)
		return
	}
	NewManagerPool(client, managers...).ManageAll(ctx)
}

// ManageAll configures every manager of the pool. Errors are logged.
func (m *ManagerPool) ManageAll(ctx context.Context) {
	for _, manager := range m.managers {
		if err := m.Configure(ctx, manager); err != nil {
			log.FromContext(ctx).Errorf("error on start vault manager. %s", err.Error())
		}
	}
}

// Configure reads the manager's secret, maps its keys and configures the
// manager. With the watcher enabled, it renews the lease in the background
// and, when the lease ends, closes and configures the manager again, until
// ctx is done.
func (m *ManagerPool) Configure(ctx context.Context, manager Manager) error {
	logger := log.FromContext(ctx)
	options := manager.Options()

	var response api.Secret
	if err := m.client.Read([]string{options.SecretPath}, &response, &vault.RequestOptions{}); err != nil {
		return err
	}

	logger.Debugf("vault secret %s read, lease %s for %vs", options.SecretPath, response.LeaseID, response.LeaseDuration)

	data := make(map[string]interface{}, len(options.Keys))
	for source, dst := range options.Keys {
		if v, ok := response.Data[source]; ok {
			data[dst] = v
		} else {
			logger.Warnf("the key %s not found in vault data", source)
		}
	}

	if err := manager.Configure(ctx, data); err != nil {
		return err
	}

	if options.Watcher.Enabled {
		go func() {
			if err := m.watch(ctx, manager, &response); err != nil {
				logger.Errorf("error on start vault watcher. %s", err.Error())
			}
		}()
	}

	return nil
}

func (m *ManagerPool) watch(ctx context.Context, manager Manager, secret *api.Secret) error {
	logger := log.FromContext(ctx)

	watcher, err := m.client.NewLifetimeWatcher(&api.LifetimeWatcherInput{
		Secret:    secret,
		Increment: manager.Options().Watcher.Increment,
	})
	if err != nil {
		return fmt.Errorf("vault: start watcher: %w", err)
	}
	go watcher.Start()
	defer watcher.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case renewal := <-watcher.RenewCh():
			logger.Debugf("vault lease renewed at %v", renewal.RenewedAt)
		case err := <-watcher.DoneCh():
			if err != nil {
				logger.Warnf("vault watcher done: %s", err.Error())
			}
			if err := manager.Close(ctx); err != nil {
				logger.Errorf("error on close vault manager. %s", err.Error())
			}
			if err := m.Configure(ctx, manager); err != nil {
				logger.Errorf("error on start vault manager. %s", err.Error())
			}
			return nil
		}
	}
}
