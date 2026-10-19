package msgpack

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	codecv2 "github.com/xgodev/boost/wrapper/cache/codec/contrib/shamaton/msgpack/v2"
)

// The shapes and values below follow wrapper/cache's manager tests, which run
// them through every codec, the v2 one included.

type complexParent struct {
	ID    string
	Name  string
	Value float64
}

type complexValue struct {
	ID        string
	CreatedAt int
	Name      string
	Parent    complexParent
}

type myStruct struct {
	ID        string
	CreatedAt int
	Name      string
}

var (
	stringValue = "value1"
	bytesValue  = []byte{0x01, 0x02, 0x03}
	structValue = complexValue{
		ID:        "123",
		CreatedAt: 2,
		Name:      "123",
		Parent:    complexParent{ID: "321", Name: "321", Value: 1.1},
	}
	sliceValue = []myStruct{
		{ID: "123", CreatedAt: 2, Name: "123"},
		{ID: "124", CreatedAt: 3, Name: "124"},
		{ID: "125", CreatedAt: 4, Name: "125"},
	}
)

func roundTrip[T any](t *testing.T, value T) {
	t.Helper()
	c := New[T]()
	b, err := c.Encode(value)
	require.NoError(t, err)
	var got T
	require.NoError(t, c.Decode(b, &got))
	assert.Equal(t, value, got)
}

// readsV2 proves a value cached by the v2 codec is read back by v3, so a
// service moving to v3 keeps the entries already in its cache.
func readsV2[T any](t *testing.T, value T) {
	t.Helper()
	b, err := codecv2.New[T]().Encode(value)
	require.NoError(t, err)
	var got T
	require.NoError(t, New[T]().Decode(b, &got))
	assert.Equal(t, value, got)
}

// v2Reads proves a value cached by v3 is read back by the v2 codec, so
// instances still on v2 during a rollout read what v3 instances wrote.
func v2Reads[T any](t *testing.T, value T) {
	t.Helper()
	b, err := New[T]().Encode(value)
	require.NoError(t, err)
	var got T
	require.NoError(t, codecv2.New[T]().Decode(b, &got))
	assert.Equal(t, value, got)
}

func TestCodec_RoundTrip(t *testing.T) {
	t.Run("string", func(t *testing.T) { roundTrip(t, stringValue) })
	t.Run("bytes", func(t *testing.T) { roundTrip(t, bytesValue) })
	t.Run("struct", func(t *testing.T) { roundTrip(t, structValue) })
	t.Run("slice of structs", func(t *testing.T) { roundTrip(t, sliceValue) })
}

func TestCodec_ReadsV2(t *testing.T) {
	t.Run("string", func(t *testing.T) { readsV2(t, stringValue) })
	t.Run("bytes", func(t *testing.T) { readsV2(t, bytesValue) })
	t.Run("struct", func(t *testing.T) { readsV2(t, structValue) })
	t.Run("slice of structs", func(t *testing.T) { readsV2(t, sliceValue) })
}

func TestCodec_V2Reads(t *testing.T) {
	t.Run("string", func(t *testing.T) { v2Reads(t, stringValue) })
	t.Run("bytes", func(t *testing.T) { v2Reads(t, bytesValue) })
	t.Run("struct", func(t *testing.T) { v2Reads(t, structValue) })
	t.Run("slice of structs", func(t *testing.T) { v2Reads(t, sliceValue) })
}

// timeValue carries a zone that is neither UTC nor the host's, so the
// assertions below hold on any machine.
var timeValue = time.Date(2026, 10, 18, 12, 30, 0, 123, time.FixedZone("BRT", -3*60*60))

// v3 decodes a MessagePack timestamp into UTC, where v2 decoded into Local:
// the instant is the same, only Location differs.
func TestCodec_Time_DecodesAsUTC(t *testing.T) {
	t.Run("written by v3", func(t *testing.T) {
		c := New[time.Time]()
		b, err := c.Encode(timeValue)
		require.NoError(t, err)
		var got time.Time
		require.NoError(t, c.Decode(b, &got))
		assert.True(t, timeValue.Equal(got), "want %v, got %v", timeValue, got)
		assert.Equal(t, time.UTC, got.Location())
	})
	t.Run("written by v2", func(t *testing.T) {
		b, err := codecv2.New[time.Time]().Encode(timeValue)
		require.NoError(t, err)
		var got time.Time
		require.NoError(t, New[time.Time]().Decode(b, &got))
		assert.True(t, timeValue.Equal(got), "want %v, got %v", timeValue, got)
		assert.Equal(t, time.UTC, got.Location())
	})
	t.Run("read by v2", func(t *testing.T) {
		b, err := New[time.Time]().Encode(timeValue)
		require.NoError(t, err)
		var got time.Time
		require.NoError(t, codecv2.New[time.Time]().Decode(b, &got))
		assert.True(t, timeValue.Equal(got), "want %v, got %v", timeValue, got)
	})
}

func TestCodec_Encode_Unsupported(t *testing.T) {
	b, err := New[chan int]().Encode(make(chan int))
	assert.Error(t, err)
	assert.Empty(t, b)
}

func TestCodec_Decode_Invalid(t *testing.T) {
	var got string
	// 0xc1 is the one byte the MessagePack spec never assigns.
	assert.Error(t, New[string]().Decode([]byte{0xc1}, &got))
	assert.Empty(t, got)
}
