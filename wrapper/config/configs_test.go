package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddKeepsOneEntryPerKey(t *testing.T) {
	saved := Entries()
	t.Cleanup(func() { SetEntries(saved) })
	SetEntries(nil)

	// two packages of one app, like the v1 and v2 pubsub adapters,
	// can register the same key; flag parsers refuse a key twice
	Add("app.key", "first", "first description")
	Add("app.other", 1, "other")
	Add("app.key", "second", "second description")

	entries := Entries()
	assert.Len(t, entries, 2)
	assert.Equal(t, "app.key", entries[0].Key)
	assert.Equal(t, "second", entries[0].Value)
	assert.Equal(t, "second description", entries[0].Description)
	assert.Equal(t, "app.other", entries[1].Key)
}
