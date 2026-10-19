package kafka

import (
	"context"
	"net"
	"strconv"
	"strings"

	"github.com/segmentio/kafka-go"
	"github.com/xgodev/boost/factory"
	"github.com/xgodev/boost/model/errors"
	"github.com/xgodev/boost/wrapper/log"
)

// DialerPlugin customizes the dialer of a connection (e.g. SASL or TLS).
type DialerPlugin = factory.Plugin[*kafka.Dialer]

// NewConn returns a new connection with default options.
func NewConn(ctx context.Context, plugins ...DialerPlugin) (*kafka.Conn, error) {
	o, err := NewOptions()
	if err != nil {
		return nil, err
	}
	return NewConnWithOptions(ctx, o, plugins...)
}

// NewConnWithConfigPath returns a new connection with options from config path.
func NewConnWithConfigPath(ctx context.Context, path string, plugins ...DialerPlugin) (*kafka.Conn, error) {
	o, err := NewOptionsWithPath(path)
	if err != nil {
		return nil, err
	}
	return NewConnWithOptions(ctx, o, plugins...)
}

// NewConnWithOptions dials the first broker. Type leader connects to the
// leader of Conn.Topic/Conn.Partition, partition to that partition on the
// first broker, server to the first broker alone.
func NewConnWithOptions(ctx context.Context, o *Options, plugins ...DialerPlugin) (*kafka.Conn, error) {

	if len(o.Brokers) == 0 {
		return nil, errors.NotValidf("kafka-go brokers are empty")
	}

	d := &kafka.Dialer{Timeout: o.Conn.Timeout, DualStack: true}
	if err := factory.ApplyPlugins(ctx, d, plugins...); err != nil {
		return nil, err
	}

	address := o.Brokers[0]

	var (
		c   *kafka.Conn
		err error
	)

	switch strings.ToLower(o.Conn.Type) {
	case "leader":
		c, err = d.DialLeader(ctx, o.Conn.Network, address, o.Conn.Topic, o.Conn.Partition)
	case "partition":
		// DialPartition dials Partition.Leader, not address: the first broker
		// is taken as the leader, with no metadata lookup.
		var leader kafka.Broker
		leader, err = broker(address)
		if err != nil {
			return nil, err
		}
		c, err = d.DialPartition(ctx, o.Conn.Network, address, kafka.Partition{
			Topic:  o.Conn.Topic,
			ID:     o.Conn.Partition,
			Leader: leader,
		})
	case "server":
		c, err = d.DialContext(ctx, o.Conn.Network, address)
	default:
		return nil, errors.NotValidf("kafka-go conn type %q", o.Conn.Type)
	}

	if err != nil {
		return nil, errors.Wrap(err, errors.Internalf("unable to connect to kafka %s: %s", address, err))
	}

	log.FromContext(ctx).Infof("connected to kafka %s", address)

	return c, nil
}

func broker(address string) (kafka.Broker, error) {
	host, p, err := net.SplitHostPort(address)
	if err != nil {
		return kafka.Broker{}, errors.NotValidf("kafka-go broker address %q: %s", address, err)
	}
	port, err := strconv.Atoi(p)
	if err != nil {
		return kafka.Broker{}, errors.NotValidf("kafka-go broker port %q: %s", address, err)
	}
	return kafka.Broker{Host: host, Port: port}, nil
}
