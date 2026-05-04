# Contributing to the boost Claude Code plugin

This directory ships **one Claude Code plugin (`boost`) containing one skill per boost subsystem**. Every new feature added to `github.com/xgodev/boost` MUST add or update the corresponding skill here.

## The Rule

> **Adding or changing a boost feature → adding or updating the corresponding skill in `cc-plugins/boost/skills/<name>/`.**

A "feature" is any change that:

- Adds a new directory under `bootstrap/function/adapter/contrib/`, `factory/contrib/`, `wrapper/*/driver/contrib/`, `wrapper/*/contrib/`, `extra/`, or `fx/modules/`.
- Adds a new public API (exported type, function, or constant) at any of the layer roots above.
- Changes a public API in a way that affects how consumers wire the subsystem (signature changes, constructor renames, new required options).
- Changes a default config key, env var, or middleware order baked into a subsystem.

Pure refactors, internal-only changes, and bug fixes that don't change the public contract do not require a skill change — but ARE allowed to motivate one if the bug fix invalidates skill content.

## Mapping: boost subsystem → skill

| Boost path | Skill (under `cc-plugins/boost/skills/`) |
|---|---|
| `bootstrap/function/` | `boost-bootstrap-function` |
| `bootstrap/function/adapter/contrib/cloud.google.com/pubsub/v1/` | `boost-bootstrap-adapter-pubsub` |
| `bootstrap/function/adapter/contrib/nats-io/nats.go/v1/` | `boost-bootstrap-adapter-nats` |
| `bootstrap/function/adapter/contrib/confluentinc/confluent-kafka-go/v2/` | `boost-bootstrap-adapter-kafka` |
| `bootstrap/function/middleware/` | `boost-bootstrap-middleware` |
| `factory/contrib/labstack/echo/v4/` | `boost-factory-echo` |
| `factory/contrib/go-resty/resty/v2/` | `boost-factory-resty` |
| `factory/contrib/cloud.google.com/pubsub/v1/` | `boost-factory-pubsub` |
| `wrapper/publisher/` | `boost-wrapper-publisher` |
| `wrapper/cache/` | `boost-wrapper-cache` |
| `wrapper/log/` | `boost-wrapper-log` |
| `wrapper/config/` | `boost-wrapper-config` |
| `fx/modules/` | `boost-fx-modules` |
| `extra/middleware/` | `boost-extra-middleware` |
| `extra/health/` | `boost-extra-health` |
| `extra/multiserver/` | `boost-extra-multiserver` |
| `boost.Start`, `wrapper/log` (basics), `wrapper/config` (basics), `model/errors` | `boost-core` |

If a feature spans multiple subsystems (e.g., a new pub/sub stack involves a new factory client + a new bootstrap adapter + a new wrapper driver), update each of the relevant skills.

If a feature touches a path NOT in the table above, add a new skill under `cc-plugins/boost/skills/<new-name>/` and add a row to this table in the same PR.

## Workflow when adding a feature

1. **Implement the feature** in the boost source tree.
2. **Identify the affected skill(s)** from the table above.
3. **For each affected skill, decide:**
   - **Lapidating a stub?** Run a TDD cycle (see below).
   - **Updating a mature skill?** Run at least one RED scenario without your change, then verify GREEN with your change. Bump the patch version in the skill's frontmatter `metadata.version`.
4. **Bump the plugin version** in `cc-plugins/boost/.claude-plugin/plugin.json` and the `marketplace.json` entry — patch bump for any skill change, minor for a new skill, major for a backwards-incompatible reorganization.
5. **Add a `CHANGELOG.md` entry** in `cc-plugins/boost/CHANGELOG.md`.
6. **Open the PR** — the `skills-coverage` GitHub Action will validate that touched feature paths have a corresponding skill change.

## TDD cycle for a stub skill (lapidation)

A stub skill's SKILL.md is a placeholder. To lapidate (turn it into a mature skill):

### RED — establish baseline

Write a scenario file under `cc-plugins/boost/skills/<name>/evals/evals.json` describing what an AI agent should produce for this subsystem (see `cc-plugins/boost/skills/boost-core/evals/evals.json` for shape). Then dispatch a fresh agent with the prompt and **without** the skill loaded — capture verbatim what they produce. Document failures and rationalizations.

### GREEN — write minimal skill

Replace the stub `SKILL.md` content with a focused skill that addresses the specific failures observed in RED. Keep the frontmatter format the same, change `metadata.status` to `mature`, and bump `metadata.version` to `0.1.0`. Re-run the same scenario; verify the agent now complies.

### REFACTOR — close loopholes

If the GREEN run produces NEW rationalizations (e.g., "but the framework forces value, not pointer"), update the skill to address them and re-run. Iterate until clean. Bump version on each iteration.

For full discipline, see the upstream `superpowers:writing-skills` skill.

## Skill file layout (for reference)

```
cc-plugins/boost/skills/<skill-name>/
├── SKILL.md
├── evals/        (after lapidation)
│   ├── evals.json
│   ├── baseline/   (gitignored — RED runs)
│   ├── green/      (gitignored — GREEN v1 runs)
│   └── green-v2/   (gitignored — GREEN v2+ runs)
└── references/   (optional)
```

## Skill versioning (in frontmatter)

Each skill uses semver in its `metadata.version`:

- `0.0.x` — stub (`metadata.status: stub`)
- `0.x.0` — lapidated, in iteration (`metadata.status: mature`)
- `1.0.0` — covered all known scenarios, no open loopholes from REFACTOR cycles

The plugin's own version (in `cc-plugins/boost/.claude-plugin/plugin.json` and `marketplace.json`) is the umbrella — bumped on every skill change so devs running `/plugin update boost@xgodev-boost` actually pull the change.

## License

The plugin is MIT, originally seeded from `samber/cc-skills-golang` (also MIT). Preserve the LICENSE file at the plugin root.
