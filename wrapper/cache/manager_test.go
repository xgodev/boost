package cache

import (
	"context"
	cfreecache "github.com/xgodev/boost/wrapper/cache/driver/contrib/coocood/freecache/v1"
	"testing"
	"time"

	"github.com/coocood/freecache"
	"github.com/stretchr/testify/suite"
	"github.com/xgodev/boost/wrapper/cache/codec/gob"
)

type ManagerSuite struct {
	suite.Suite
}

func TestManagerSuite(t *testing.T) {
	suite.Run(t, new(ManagerSuite))
}

func (s *ManagerSuite) SetupSuite() {}

func (s *ManagerSuite) TearDownSuite() {}

func (s *ManagerSuite) Test_manager_Del() {

	tt := []struct {
		name    string
		key     string
		data    map[string]string
		wantErr bool
	}{
		{
			name: "when exists key",
			key:  "key1",
			data: map[string]string{
				"key1": "value1",
			},
			wantErr: false,
		},
		{
			name: "when not exists key",
			key:  "key2",
			data: map[string]string{
				"key1": "value1",
			},
			wantErr: false,
		},
	}

	for _, t := range tt {
		s.Run(t.name, func() {

			ctx := context.Background()

			fc := freecache.NewCache(1)
			drv := cfreecache.New(fc, &cfreecache.Options{TTL: 1 * time.Minute})

			codec := gob.New[string]()

			manager := NewManager[string]("foo", codec, drv)
			for key, value := range t.data {
				s.Assert().Nil(manager.Set(ctx, key, value))
			}

			err := manager.Del(ctx, t.key)
			var got bool
			if err != nil {
				got = true
			}
			s.Assert().True(got == t.wantErr)
		})
	}
}

func (s *ManagerSuite) Test_manager_Get() {

	tt := []struct {
		name    string
		key     string
		data1   map[string]string
		data2   map[string]string
		wantErr bool
		want    string
	}{
		{
			name: "when exists key",
			key:  "key1",
			data1: map[string]string{
				"key1": "value1",
			},
			data2: map[string]string{
				"key1": "value1",
			},
			want:    "value1",
			wantErr: false,
		},
		{
			name: "when key exists on second level",
			key:  "key2",
			data1: map[string]string{
				"key1": "value1",
			},
			data2: map[string]string{
				"key2": "value2",
			},
			want:    "value2",
			wantErr: false,
		},
	}

	for _, t := range tt {
		s.Run(t.name, func() {

			ctx := context.Background()

			codec := gob.New[string]()

			fc1 := freecache.NewCache(1)
			drv1 := cfreecache.New(fc1, &cfreecache.Options{TTL: 1 * time.Minute})

			for key, value := range t.data1 {
				b, _ := codec.Encode(value)
				s.Assert().Nil(drv1.Set(ctx, key, b))
			}

			fc2 := freecache.NewCache(1)
			drv2 := cfreecache.New(fc2, &cfreecache.Options{TTL: 1 * time.Minute})

			for key, value := range t.data2 {
				b, _ := codec.Encode(value)
				s.Assert().Nil(drv2.Set(ctx, key, b))
			}

			manager := NewManager[string]("foo", codec, drv1, drv2)
			for key, value := range t.data1 {
				s.Assert().Nil(manager.Set(ctx, key, value))
			}

			_, got, err := manager.Get(ctx, t.key)
			var gotErr bool
			if err != nil {
				gotErr = true
			}
			s.Assert().True(gotErr == t.wantErr)
			s.Assert().Equal(t.want, got)
		})
	}
}

/*
func (s *ManagerSuite) Test_manager_Set() {

	tt := []struct {
		name    string
		key     string
		data    map[string]string
		wantErr bool
	}{
		{
			name: "when exists key",
			key:  "key1",
			data: map[string]string{
				"key1": "value1",
			},
			wantErr: false,
		},
		{
			name: "when not exists key",
			key:  "key2",
			data: map[string]string{
				"key1": "value1",
			},
			wantErr: false,
		},
	}

	for _, t := range tt {
		s.Run(t.name, func() {

			ctx := context.Background()

			cfg := bigcache.DefaultConfig(1 * time.Minute)
			cfg.HardMaxCacheSize = 1
			cfg.Shards = 2
			fc, _ := bigcache.New(ctx, cfg)
			drv := New(fc)

			codec := gob.New[string]()

			for key, value := range t.data {
				b, _ := codec.Encode(value)
				s.Assert().Nil(drv.Set(ctx, key, b))
			}

			v := "xpto"

			b, _ := codec.Encode(v)
			err := drv.Set(ctx, t.key, b)
			var gotErr bool
			if err != nil {
				gotErr = true
			}
			s.Assert().True(gotErr == t.wantErr)

			var got string
			b, err = drv.Get(ctx, t.key)
			codec.Decode(b, &got)

			s.Assert().Equal(v, got)
		})
	}
}

*/
