# Changelog

All notable changes to the `boost` Claude Code plugin.

The plugin uses semver. Each individual skill carries its own `metadata.version` in the SKILL.md frontmatter — the plugin version bumps when any skill's content or metadata changes.

## [0.3.0] — initial bundle

- Bundled `boost-core` (mature, v0.3.0) — Iron Laws, `boost.Start`, `log.FromContext`, `config`, `model/errors`.
- Bundled 17 subsystem stubs (each at v0.0.1, status `stub`) covering `factory/`, `bootstrap/`, `wrapper/`, `fx/`, `extra/`, and the `boost-maintainer` contributor guide.
- Each stub ships with a lapidation checklist; converted to mature via TDD cycle when first exercised by a real boost feature change.
