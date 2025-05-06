package cache

import (
	"context"

	"github.com/xgodev/boost/wrapper/log"
)

type Manager[T any] struct {
	drivers []Driver
	mids    []Plugin[T]
	codec   Codec[T]
	metric  *Metric
	name    string
}

func (m *Manager[T]) newContext(ctx context.Context, driver Driver) *Context[T] {
	c := NewContext[T](m.name, driver, m.mids...)
	c.SetContext(ctx)
	return c
}

func (m *Manager[T]) Del(ctx context.Context, key string) error {

	for _, d := range m.drivers {
		c := m.newContext(ctx, d)
		if err := c.Del(key); err != nil {
			return err
		}
	}
	return nil
}

func (m *Manager[T]) Get(ctx context.Context, key string) (ok bool, data T, err error) {

	var b []byte

	for _, d := range m.drivers {
		c := m.newContext(ctx, d)
		b, err = c.Get(key)
		if err != nil {
			return false, data, err
		}
		if len(b) > 0 {
			break
		}
	}

	if len(b) > 0 {
		d := data
		if err = m.codec.Decode(b, &d); err != nil {
			return false, data, err
		}
		return true, d, err
	}

	return false, data, err
}

func (m *Manager[T]) Set(ctx context.Context, key string, data T, opts ...OptionSet) (err error) {
	var b []byte

	if b, err = m.codec.Encode(data); err != nil {
		return err
	}

	return m.set(ctx, len(m.drivers)-1, key, b, opts...)
}

func (m *Manager[T]) set(ctx context.Context, driveIndex int, key string, b []byte, opts ...OptionSet) (err error) {

	opt := Option{
		SaveEmpty: false,
		AsyncSave: false,
		Replicate: true,
	}

	for _, o := range opts {
		o()(&opt)
	}

	if len(b) > 0 || opt.SaveEmpty {

		if opt.AsyncSave {

			go func(ctx context.Context, key string, b []byte) {
				for i, d := range m.drivers {
					if m.isReplicable(opt, i, driveIndex) {
						c := m.newContext(ctx, d)
						if err := c.Set(key, b); err != nil {
							log.Error(err.Error())
						}
					}
				}
			}(ctx, key, b)

		} else {

			for i, d := range m.drivers {
				if m.isReplicable(opt, i, driveIndex) {
					c := m.newContext(ctx, d)
					if err = c.Set(key, b); err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}

func (m *Manager[T]) isReplicable(opt Option, i int, driveIndex int) bool {
	return (opt.Replicate && i < driveIndex) || len(m.drivers)-1 == driveIndex
}

func (m *Manager[T]) GetOrSet(ctx context.Context, key string, cacheable Cacheable[T], opts ...OptionSet) (data T, err error) {

	var b []byte
	var index int

	for i, d := range m.drivers {
		c := m.newContext(ctx, d)
		b, err = c.Get(key)
		if err != nil {
			return data, err
		}
		if len(b) > 0 {
			index = i
			break
		}
	}

	if len(b) > 0 {
		if err = m.codec.Decode(b, &data); err != nil {
			return data, err
		}

		if index > 0 {
			err = m.set(ctx, index, key, b, opts...)
			if err != nil {
				return data, err
			}
		}

	} else {
		data, err = cacheable(ctx)
		if err != nil {
			return data, err
		}

		err = m.Set(ctx, key, data, opts...)
		if err != nil {
			return data, err
		}
	}

	return data, err
}

func (m *Manager[T]) Use(mid Plugin[T]) *Manager[T] {
	m.mids = append(m.mids, mid)
	return m
}

func NewManager[T any](name string, c Codec[T], d ...Driver) *Manager[T] {
	return &Manager[T]{name: name, codec: c, drivers: d}
}
