package freecache

import (
	"context"
	"testing"
	"time"

	"github.com/coocood/freecache"
	"github.com/stretchr/testify/suite"
	"github.com/xgodev/boost/wrapper/cache/codec/gob"
)

type DriverSuite struct {
	suite.Suite
}

func TestConfigSuite(t *testing.T) {
	suite.Run(t, new(DriverSuite))
}

func (s *DriverSuite) SetupSuite() {}

func (s *DriverSuite) TearDownSuite() {}

func (s *DriverSuite) Test_driver_Del() {

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
			drv := New(fc, &Options{TTL: 1 * time.Minute})

			codec := gob.New[string]()

			for key, value := range t.data {
				b, _ := codec.Encode(value)
				s.Assert().Nil(drv.Set(ctx, key, b))
			}

			err := drv.Del(ctx, t.key)
			var got bool
			if err != nil {
				got = true
			}
			s.Assert().True(got == t.wantErr)
		})
	}
}

func (s *DriverSuite) Test_driver_Get() {

	tt := []struct {
		name    string
		key     string
		data    map[string]string
		wantErr bool
		want    string
	}{
		{
			name: "when exists key",
			key:  "key1",
			data: map[string]string{
				"key1": "value1",
			},
			want:    "value1",
			wantErr: false,
		},
		{
			name: "when not exists key",
			key:  "key2",
			data: map[string]string{
				"key1": "value1",
			},
			want:    "",
			wantErr: false,
		},
	}

	for _, t := range tt {
		s.Run(t.name, func() {

			ctx := context.Background()

			fc := freecache.NewCache(1)
			drv := New(fc, &Options{TTL: 1 * time.Minute})

			codec := gob.New[string]()

			for key, value := range t.data {
				b, _ := codec.Encode(value)
				s.Assert().Nil(drv.Set(ctx, key, b))
			}

			var got string
			b, err := drv.Get(ctx, t.key)
			codec.Decode(b, &got)

			var gotErr bool
			if err != nil {
				gotErr = true
			}
			s.Assert().True(gotErr == t.wantErr)
			s.Assert().Equal(t.want, got)
		})
	}
}

func (s *DriverSuite) Test_driver_Set() {

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
			drv := New(fc, &Options{TTL: 1 * time.Minute})

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
