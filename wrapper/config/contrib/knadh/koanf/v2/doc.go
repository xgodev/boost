// Package koanf is the boost config driver on koanf v2.
//
// Load reads, in this order, the map defaults declared with config.Add, the
// files listed in the CONF environment variable or the --conf flag (.json,
// .yaml, .yml, .toml), the environment and the command line; each source
// overrides the previous one.
//
// The core (github.com/knadh/koanf/v2), the env provider (providers/env/v2)
// and the TOML parser (parsers/toml/v2) are koanf's own modules. The flag,
// file, JSON and YAML adapters live in this package: koanf's modules for them
// share their import paths with the github.com/knadh/koanf v1 module, which
// boost still requires for the v1 driver, and Go refuses both in one build.
package koanf
