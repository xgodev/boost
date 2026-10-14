// Package names contributes values to the "names" group.
package names

// Alice adds Alice to the names group.
// @Provide(group=names, index=0)
func Alice() string {
	return "Alice"
}

// Bob adds Bob to the names group.
// @Provide(group=names, index=0)
func Bob() string {
	return "Bob"
}
