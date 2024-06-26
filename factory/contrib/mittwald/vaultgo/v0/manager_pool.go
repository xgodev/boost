package vault

import (
	"context"

	"github.com/hashicorp/vault/api"
	vault "github.com/mittwald/vaultgo"
	"github.com/xgodev/boost/model/errors"
	"github.com/xgodev/boost/wrapper/log"
)

// ManagerPool represents a manager pool for vault.
type ManagerPool struct {
	managers []Manager
	client   *vault.Client
}

// NewManagerPool returns a new manager pool.
func NewManagerPool(client *vault.Client, managers ...Manager) *ManagerPool {
	return &ManagerPool{managers: managers, client: client}
}

// ManageAll configures all managers to a new manager pool.
func ManageAll(ctx context.Context, managers ...Manager) {
	client, err := NewClient(ctx)
	if err != nil {
		log.Error(err)
		return
	}
	mp := NewManagerPool(client, managers...)
	mp.ManageAll(ctx)
}

// ManageAll configures all managers from this manager pool.
func (m *ManagerPool) ManageAll(ctx context.Context) {

	for _, manager := range m.managers {
		mr := manager
		err := m.Configure(ctx, mr)
		if err != nil {
			log.Errorf("error on start vault manager. %s", err.Error())
		}
	}

}

// Configure configures manager with options from this manager pool.
func (m *ManagerPool) Configure(ctx context.Context, manager Manager) error {

	var response api.Secret

	err := m.client.Read([]string{manager.Options().SecretPath}, &response, &vault.RequestOptions{
		Parameters:  nil,
		SkipRenewal: false,
	})
	if err != nil {
		return err
	}

	log.Debugf("lease_id: %s", response.LeaseID)
	log.Debugf("data: %v", response.Data)
	log.Debugf("lease_duration: %vs", response.LeaseDuration)

	data := response.Data
	dataConv := make(map[string]interface{})

	options := manager.Options()

	for source, dst := range options.Keys {
		if dt, ok := data[source]; ok {
			dataConv[dst] = dt
		} else {
			log.Warnf("the key %s not found in vault data", source)
		}
	}

	if err := manager.Configure(ctx, dataConv); err != nil {
		return err
	}

	if manager.Options().Watcher.Enabled {
		go func() {
			err := m.watch(ctx, manager, response)
			if err != nil {
				log.Errorf("error on start vault watcher. %s", err.Error())
			}
		}()
	}

	return nil
}

func (m *ManagerPool) watch(ctx context.Context, manager Manager, response api.Secret) error {

	secretesTokenWatcher := api.LifetimeWatcherInput{
		Secret:    &response,
		Increment: manager.Options().Watcher.Increment,
	}
	watcher, err := m.client.NewLifetimeWatcher(&secretesTokenWatcher)
	if err != nil {
		return errors.Internalf("error on start watcher. %s", err.Error())
	}
	go watcher.Start()

	for {
		select {
		case rawData := <-watcher.RenewCh():
			log.Debugf("received renewal at: %+v", rawData.RenewedAt)
			log.Debugf("received renewal Secret: %+v", rawData.Secret)
		case err := <-watcher.DoneCh():
			if err != nil {
				log.Errorf("Got watcher error: %s", err.Error())
			}
			watcher.Stop()
			if er := manager.Close(ctx); er != nil {
				log.Errorf("Got manager error: %s", er.Error())
			}
			go func() {
				err := m.Configure(ctx, manager)
				if err != nil {
					log.Errorf("error on start vault manager. %s", err.Error())
				}
			}()
			return nil
		}
	}

}
