package cache

import (
	"context"
	"fmt"
	"github.com/allegro/bigcache/v3"
	codec_binary "github.com/xgodev/boost/wrapper/cache/codec/binary"
	codec_goccy "github.com/xgodev/boost/wrapper/cache/codec/contrib/goccy/go-json/v0"
	codec_shamaton_msgpack "github.com/xgodev/boost/wrapper/cache/codec/contrib/shamaton/msgpack/v2"
	codec_vmihailenco_msgpack "github.com/xgodev/boost/wrapper/cache/codec/contrib/vmihailenco/msgpack/v5"
	codec_gob "github.com/xgodev/boost/wrapper/cache/codec/gob"
	codec_json "github.com/xgodev/boost/wrapper/cache/codec/json"
	codec_string "github.com/xgodev/boost/wrapper/cache/codec/string"
	cbigcache "github.com/xgodev/boost/wrapper/cache/driver/contrib/allegro/bigcache/v3"
	cfreecache "github.com/xgodev/boost/wrapper/cache/driver/contrib/coocood/freecache/v1"
	"testing"
	"time"

	"github.com/coocood/freecache"
	"github.com/stretchr/testify/suite"
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

			codec := codec_string.New[string]()

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

func (s *ManagerSuite) Test_manager_Get_WithReplicate() {

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

			codec := codec_gob.New[string]()

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

func (s *ManagerSuite) Test_manager_Codecs_SimpleString() {

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
	}

	fc := freecache.NewCache(1)
	drv := cfreecache.New(fc, &cfreecache.Options{TTL: 1 * time.Minute})

	for _, t := range tt {
		s.Run(t.name, func() {

			ctx := context.Background()

			codecs := map[string]Codec[string]{
				"gob":                 codec_gob.New[string](),
				"msgpack_vmihailenco": codec_vmihailenco_msgpack.New[string](),
				"msgpack_shamaton":    codec_shamaton_msgpack.New[string](),
				"string":              codec_string.New[string](),
				"goccy":               codec_goccy.New[string](),
				"json":                codec_json.New[string](),
			}

			// Variables to store the best results
			var fastestCodec string
			var smallestCodec string
			var minTime time.Duration = time.Hour // Set to a high value initially
			var minSize int = int(^uint(0) >> 1)  // Set to the maximum possible int value

			// Map to store the time and size results of each codec
			results := make(map[string]struct {
				Time time.Duration
				Size int
			})

			for key, value := range t.data {
				for name, codec := range codecs {
					start := time.Now()

					// Encode the value
					b, err := codec.Encode(value)
					duration := time.Since(start)

					if err != nil {
						s.Assert().Failf(err.Error(), "error on codec: %s", name)
					}

					// Measure the size of the encoded result
					size := len(b)

					// Store the results for each codec
					results[name] = struct {
						Time time.Duration
						Size int
					}{
						Time: duration,
						Size: size,
					}

					// Track the fastest codec
					if duration < minTime {
						minTime = duration
						fastestCodec = name
					}

					// Track the codec with the smallest size
					if size < minSize {
						minSize = size
						smallestCodec = name
					}

					s.Assert().Nilf(drv.Set(ctx, key, b), "codec: %s", name)
					res, err := drv.Get(ctx, key)
					s.Assert().Nilf(err, "codec: %s", name)
					s.Assert().NotEmpty(res, "codec: %s", name)
				}
			}

			// Print results for all codecs
			fmt.Println("Codec Performance Results:")
			for name, result := range results {
				fmt.Printf("Codec: %s - Time: %v - Size: %d bytes\n", name, result.Time, result.Size)
			}

			// Print the fastest and most compact results
			fmt.Printf("\nFastest codec: %s (%v)\n", fastestCodec, minTime)
			fmt.Printf("Codec with smallest size: %s (%d bytes)\n", smallestCodec, minSize)

			for name, codec := range codecs {
				manager := NewManager[string]("foo", codec, drv)
				for key, value := range t.data {
					s.Assert().Nilf(manager.Set(ctx, key, value), "codec: %s", name)
					res, err := drv.Get(ctx, key)
					s.Assert().Nilf(err, "codec: %s", name)
					s.Assert().NotEmpty(res, "codec: %s", name)
				}

				_, got, err := manager.Get(ctx, t.key)
				var gotErr bool
				if err != nil {
					gotErr = true
				}
				s.Assert().True(gotErr == t.wantErr)
				s.Assert().Equal(t.want, got)
			}
		})
	}
}

func (s *ManagerSuite) Test_manager_Codecs_ByteArray() {

	tt := []struct {
		name    string
		key     string
		data    map[string][]byte
		wantErr bool
		want    []byte
	}{
		{
			name: "when exists key",
			key:  "key1",
			data: map[string][]byte{
				"key1": []byte{0x01, 0x02, 0x03},
			},
			want:    []byte{0x01, 0x02, 0x03},
			wantErr: false,
		},
	}

	fc := freecache.NewCache(1)
	drv := cfreecache.New(fc, &cfreecache.Options{TTL: 1 * time.Minute})

	for _, t := range tt {
		s.Run(t.name, func() {

			ctx := context.Background()

			// Defining codecs for []byte
			codecs := map[string]Codec[[]byte]{
				"binary":              codec_binary.New[[]byte](),
				"gob":                 codec_gob.New[[]byte](),
				"msgpack_vmihailenco": codec_vmihailenco_msgpack.New[[]byte](),
				"msgpack_shamaton":    codec_shamaton_msgpack.New[[]byte](),
				"goccy":               codec_goccy.New[[]byte](),
				"json":                codec_json.New[[]byte](),
			}

			// Variables to store the best results
			var fastestCodec string
			var smallestCodec string
			var minTime time.Duration = time.Hour // Set to a high value initially
			var minSize int = int(^uint(0) >> 1)  // Set to the maximum possible int value

			// Map to store the time and size results of each codec
			results := make(map[string]struct {
				Time time.Duration
				Size int
			})

			for key, value := range t.data {
				for name, codec := range codecs {
					start := time.Now()

					// Encode the value
					b, err := codec.Encode(value)
					duration := time.Since(start)

					if err != nil {
						s.Assert().Failf(err.Error(), "error on codec: %s", name)
					}

					// Measure the size of the encoded result
					size := len(b)

					// Store the results for each codec
					results[name] = struct {
						Time time.Duration
						Size int
					}{
						Time: duration,
						Size: size,
					}

					// Track the fastest codec
					if duration < minTime {
						minTime = duration
						fastestCodec = name
					}

					// Track the codec with the smallest size
					if size < minSize {
						minSize = size
						smallestCodec = name
					}

					s.Assert().Nilf(drv.Set(ctx, key, b), "codec: %s", name)
					res, err := drv.Get(ctx, key)
					s.Assert().Nilf(err, "codec: %s", name)
					s.Assert().NotEmpty(res, "codec: %s", name)
				}
			}

			// Print results for all codecs
			fmt.Println("Codec Performance Results:")
			for name, result := range results {
				fmt.Printf("Codec: %s - Time: %v - Size: %d bytes\n", name, result.Time, result.Size)
			}

			// Print the fastest and most compact results
			fmt.Printf("\nFastest codec: %s (%v)\n", fastestCodec, minTime)
			fmt.Printf("Codec with smallest size: %s (%d bytes)\n", smallestCodec, minSize)

			for name, codec := range codecs {
				manager := NewManager[[]byte]("foo", codec, drv)
				for key, value := range t.data {
					s.Assert().Nilf(manager.Set(ctx, key, value), "codec: %s", name)
					res, err := drv.Get(ctx, key)
					s.Assert().Nilf(err, "codec: %s", name)
					s.Assert().NotEmpty(res, "codec: %s", name)
				}

				_, got, err := manager.Get(ctx, t.key)
				var gotErr bool
				if err != nil {
					gotErr = true
				}
				s.Assert().True(gotErr == t.wantErr)
				s.Assert().Equal(t.want, got)
			}
		})
	}
}

func (s *ManagerSuite) Test_manager_Codecs_Struct() {

	type ComplexParent struct {
		ID    string
		Name  string
		Value float64
	}

	type Complex struct {
		ID        string
		CreatedAt int
		Name      string
		Parent    ComplexParent
	}

	tt := []struct {
		name    string
		key     string
		data    map[string]Complex
		wantErr bool
		want    Complex
	}{
		{
			name: "when exists key",
			key:  "key1",
			data: map[string]Complex{
				"key1": {
					ID:        "123",
					CreatedAt: 2,
					Name:      "123",
					Parent: ComplexParent{
						ID:    "321",
						Name:  "321",
						Value: 1.1,
					},
				},
			},
			want: Complex{
				ID:        "123",
				CreatedAt: 2,
				Name:      "123",
				Parent: ComplexParent{
					ID:    "321",
					Name:  "321",
					Value: 1.1,
				},
			},
			wantErr: false,
		},
	}

	fc := freecache.NewCache(1)
	drv := cfreecache.New(fc, &cfreecache.Options{TTL: 1 * time.Minute})

	for _, t := range tt {
		s.Run(t.name, func() {

			ctx := context.Background()

			codecs := map[string]Codec[Complex]{
				"gob":                 codec_gob.New[Complex](),
				"msgpack_vmihailenco": codec_vmihailenco_msgpack.New[Complex](),
				"msgpack_shamaton":    codec_shamaton_msgpack.New[Complex](),
				"json":                codec_json.New[Complex](),
				"goccy":               codec_goccy.New[Complex](),
			}

			// Variables to store the best results
			var fastestCodec string
			var smallestCodec string
			var minTime time.Duration = time.Hour // Set to a high value initially
			var minSize int = int(^uint(0) >> 1)  // Set to the maximum possible int value

			// Map to store the time and size results of each codec
			results := make(map[string]struct {
				Time time.Duration
				Size int
			})

			for key, value := range t.data {
				for name, codec := range codecs {
					start := time.Now()

					// Encode the value
					b, err := codec.Encode(value)
					duration := time.Since(start)

					if err != nil {
						s.Assert().Failf(err.Error(), "error on encode codec: %s", name)
					}

					// Measure the size of the encoded result
					size := len(b)

					// Store the results for each codec
					results[name] = struct {
						Time time.Duration
						Size int
					}{
						Time: duration,
						Size: size,
					}

					// Track the fastest codec
					if duration < minTime {
						minTime = duration
						fastestCodec = name
					}

					// Track the codec with the smallest size
					if size < minSize {
						minSize = size
						smallestCodec = name
					}

					s.Assert().Nilf(drv.Set(ctx, key, b), "codec: %s", name)
					res, err := drv.Get(ctx, key)
					s.Assert().Nilf(err, "codec: %s", name)
					s.Assert().NotEmpty(res, "codec: %s", name)
				}
			}

			// Print results for all codecs
			fmt.Println("Codec Performance Results:")
			for name, result := range results {
				fmt.Printf("Codec: %s - Time: %v - Size: %d bytes\n", name, result.Time, result.Size)
			}

			// Print the fastest and most compact results
			fmt.Printf("\nFastest codec: %s (%v)\n", fastestCodec, minTime)
			fmt.Printf("Codec with smallest size: %s (%d bytes)\n", smallestCodec, minSize)

			for name, codec := range codecs {
				manager := NewManager[Complex]("foo", codec, drv)
				for key, value := range t.data {
					s.Assert().Nilf(manager.Set(ctx, key, value), "codec: %s", name)
					res, err := drv.Get(ctx, key)
					s.Assert().Nilf(err, "codec: %s", name)
					s.Assert().NotEmpty(res, "codec: %s", name)
				}

				_, got, err := manager.Get(ctx, t.key)
				var gotErr bool
				if err != nil {
					gotErr = true
				}
				s.Assert().True(gotErr == t.wantErr)
				s.Assert().Equal(t.want, got)
			}
		})
	}
}

func (s *ManagerSuite) Test_manager_Codecs_BigData() {

	type MyStruct struct {
		ID        string
		CreatedAt int
		Name      string
	}

	tt := []struct {
		name string
		key  string
		data map[string][]MyStruct
	}{
		{
			name: "when exists key",
			key:  "key1",
			data: map[string][]MyStruct{
				"key1": {
					{
						ID:        "123",
						CreatedAt: 2,
						Name:      "123",
					},
					{
						ID:        "123",
						CreatedAt: 2,
						Name:      "123",
					},
					{
						ID:        "123",
						CreatedAt: 2,
						Name:      "123",
					},
					{
						ID:        "123",
						CreatedAt: 2,
						Name:      "123",
					},
					{
						ID:        "123",
						CreatedAt: 2,
						Name:      "123",
					},
					{
						ID:        "123",
						CreatedAt: 2,
						Name:      "123",
					},
					{
						ID:        "123",
						CreatedAt: 2,
						Name:      "123",
					},
					{
						ID:        "123",
						CreatedAt: 2,
						Name:      "123",
					},
					{
						ID:        "123",
						CreatedAt: 2,
						Name:      "123",
					},
					{
						ID:        "123",
						CreatedAt: 2,
						Name:      "123",
					},
					{
						ID:        "123",
						CreatedAt: 2,
						Name:      "123",
					},
					{
						ID:        "123",
						CreatedAt: 2,
						Name:      "123",
					},
					{
						ID:        "123",
						CreatedAt: 2,
						Name:      "123",
					},
					{
						ID:        "123",
						CreatedAt: 2,
						Name:      "123",
					},
					{
						ID:        "123",
						CreatedAt: 2,
						Name:      "123",
					},
					{
						ID:        "123",
						CreatedAt: 2,
						Name:      "123",
					},
					{
						ID:        "123",
						CreatedAt: 2,
						Name:      "123",
					},
					{
						ID:        "123",
						CreatedAt: 2,
						Name:      "123",
					},
					{
						ID:        "123",
						CreatedAt: 2,
						Name:      "123",
					},
					{
						ID:        "123",
						CreatedAt: 2,
						Name:      "123",
					},
					{
						ID:        "123",
						CreatedAt: 2,
						Name:      "123",
					},
					{
						ID:        "123",
						CreatedAt: 2,
						Name:      "123",
					},
					{
						ID:        "123",
						CreatedAt: 2,
						Name:      "123",
					},
					{
						ID:        "123",
						CreatedAt: 2,
						Name:      "123",
					},
					{
						ID:        "123",
						CreatedAt: 2,
						Name:      "123",
					},
					{
						ID:        "123",
						CreatedAt: 2,
						Name:      "123",
					},
					{
						ID:        "123",
						CreatedAt: 2,
						Name:      "123",
					},
					{
						ID:        "123",
						CreatedAt: 2,
						Name:      "123",
					},
					{
						ID:        "123",
						CreatedAt: 2,
						Name:      "123",
					},
					{
						ID:        "123",
						CreatedAt: 2,
						Name:      "123",
					},
					{
						ID:        "123",
						CreatedAt: 2,
						Name:      "123",
					},
					{
						ID:        "123",
						CreatedAt: 2,
						Name:      "123",
					},
					{
						ID:        "123",
						CreatedAt: 2,
						Name:      "123",
					},
					{
						ID:        "123",
						CreatedAt: 2,
						Name:      "123",
					},
					{
						ID:        "123",
						CreatedAt: 2,
						Name:      "123",
					},
					{
						ID:        "123",
						CreatedAt: 2,
						Name:      "123",
					},
					{
						ID:        "123",
						CreatedAt: 2,
						Name:      "123",
					},
					{
						ID:        "123",
						CreatedAt: 2,
						Name:      "123",
					},
					{
						ID:        "123",
						CreatedAt: 2,
						Name:      "123",
					},
					{
						ID:        "123",
						CreatedAt: 2,
						Name:      "123",
					},
					{
						ID:        "123",
						CreatedAt: 2,
						Name:      "123",
					},
					{
						ID:        "123",
						CreatedAt: 2,
						Name:      "123",
					},
					{
						ID:        "123",
						CreatedAt: 2,
						Name:      "123",
					},
					{
						ID:        "123",
						CreatedAt: 2,
						Name:      "123",
					},
					{
						ID:        "123",
						CreatedAt: 2,
						Name:      "123",
					},
					{
						ID:        "123",
						CreatedAt: 2,
						Name:      "123",
					},
				},
			},
		},
	}

	fc := freecache.NewCache(10 * 1024 * 1024)
	drv := cfreecache.New(fc, &cfreecache.Options{TTL: 1 * time.Minute})

	for _, t := range tt {
		s.Run(t.name, func() {

			ctx := context.Background()

			codecs := map[string]Codec[[]MyStruct]{
				"gob":                 codec_gob.New[[]MyStruct](),
				"msgpack_vmihailenco": codec_vmihailenco_msgpack.New[[]MyStruct](),
				"msgpack_shamaton":    codec_shamaton_msgpack.New[[]MyStruct](),
				"json":                codec_json.New[[]MyStruct](),
				"goccy":               codec_goccy.New[[]MyStruct](),
			}

			// Variables to store the best results
			var fastestCodec string
			var smallestCodec string
			var minTime time.Duration = time.Hour // Set to a high value initially
			var minSize int = int(^uint(0) >> 1)  // Set to the maximum possible int value

			// Map to store the time and size results of each codec
			results := make(map[string]struct {
				Time time.Duration
				Size int
			})

			for key, value := range t.data {
				for name, codec := range codecs {
					start := time.Now()

					// Encode the value
					b, err := codec.Encode(value)
					duration := time.Since(start)

					if err != nil {
						s.Assert().Failf(err.Error(), "error on encode codec: %s", name)
					}

					// Measure the size of the encoded result
					size := len(b)

					// Store the results for each codec
					results[name] = struct {
						Time time.Duration
						Size int
					}{
						Time: duration,
						Size: size,
					}

					// Track the fastest codec
					if duration < minTime {
						minTime = duration
						fastestCodec = name
					}

					// Track the codec with the smallest size
					if size < minSize {
						minSize = size
						smallestCodec = name
					}

					s.Assert().Nilf(drv.Set(ctx, key, b), "codec: %s", name)
					res, err := drv.Get(ctx, key)
					s.Assert().Nilf(err, "codec: %s", name)
					s.Assert().NotEmpty(res, "codec: %s", name)
				}
			}

			// Print results for all codecs
			fmt.Println("Codec Performance Results:")
			for name, result := range results {
				fmt.Printf("Codec: %s - Time: %v - Size: %d bytes\n", name, result.Time, result.Size)
			}

			// Print the fastest and most compact results
			fmt.Printf("\nFastest codec: %s (%v)\n", fastestCodec, minTime)
			fmt.Printf("Codec with smallest size: %s (%d bytes)\n", smallestCodec, minSize)

			for name, codec := range codecs {
				manager := NewManager[[]MyStruct]("foo", codec, drv)
				for key, value := range t.data {
					s.Assert().Nilf(manager.Set(ctx, key, value), "codec: %s", name)
					res, err := drv.Get(ctx, key)
					s.Assert().Nilf(err, "codec: %s", name)
					s.Assert().NotEmpty(res, "codec: %s", name)
				}
			}
		})
	}
}

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
			drv := cbigcache.New(fc)

			codec := codec_gob.New[string]()

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
