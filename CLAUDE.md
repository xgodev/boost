# boost — Claude Code

Framework modular e extensível para serviços Go (`github.com/xgodev/boost`):
boot único (`boost.Start`), config/log/cache/publisher como wrappers, e uma
factory por componente sob `factory/contrib/`.

## Iron Laws — NUNCA violar

1. **`boost.Start()` primeiro.** Config e logger só existem depois dele.
   Nada de `config`/`log` antes do `Start`.
2. **Log só via `wrapper/log`.** `log.FromContext(ctx)` — nunca
   `zap.NewProduction()`, `zerolog.New()`, `logrus.New()` em código de app.
3. **Config só via `wrapper/config`.** Nada de `os.Getenv`/viper direto fora
   do `config.Add` em `init()`. Namespacing `boost.factory.<x>.*`.
4. **Uma factory por componente, via construtor da factory.** Nunca instanciar
   o SDK upstream direto — use `factory/contrib/<x>`.
5. **Compatibilidade retroativa.** API pública de um componente não quebra sem
   bump de major e justificativa. Mudança observável → teste antes.
6. **O plugin `golang-boost` vive em `xgodev/boost-claude` e co-evolui.** A
   doc de consumo lá é uma única skill (`skills/boost/`) com folhas em
   `skills/boost/references/<grupo>/<nome>.md`. Mudou um componente → siga a
   skill de projeto **`boost-claude-sync`** (`.claude/skills/`): atualizar a
   reference correspondente, índices e bump do `plugin.json`, em PR próprio
   lá. Sem bump, auto-update não reconhece. (O acoplamento "mesmo PR" virou
   disciplina entre dois repos; sem o PR no boost-claude, a doc fica
   defasada.)

### Red flags — PARAR e reportar

- SDK upstream instanciado direto (sem passar pela factory)
- `os.Getenv` / logger de terceiro fora do wrapper
- API pública alterada sem teste e sem nota de compat
- Componente tocado sem a reference correspondente atualizada em `xgodev/boost-claude` (skill `boost-claude-sync` não seguida)
- `git push` sem quality-gate verde
- Issue ou PR ausentes (ver `docs/development/gitflow.md`)

## Regras gerais de código

- **`gofmt`/`goimports` limpo. `go vet` e linter sem warning.**
- **`go build ./...` e `go test ./...` verdes** antes de qualquer push.
- **Single source of truth** — config keys e defaults declarados uma vez no
  `config.go` do componente.
- **Separação de concerns** — factory constrói, wrapper abstrai, app consome.
- Documentação é parte da tarefa: mudou componente/config key/comportamento →
  atualizar README do pacote **e** a skill `boost-*` no mesmo commit.
- Teste valida comportamento, não só cobre linha. Bug fix → teste vermelho
  primeiro (TDD).

## Doc do plugin (vive em `xgodev/boost-claude`)

A doc de consumo do boost vive no repo
**[`xgodev/boost-claude`](https://github.com/xgodev/boost-claude)** como uma
única skill de entrada `skills/boost/` (índice + folhas em
`skills/boost/references/<grupo>/<nome>.md`). Este repo (`boost`) **não é
marketplace nem contém skills do plugin**.

- `skills/boost/references/<grupo>/<nome>.md` — doc de consumo por
  componente (como usar boost).
- `skills/boost/references/CONTRIBUTING.md` — guia de manutenção (layout,
  construtores, config). Leia antes de adicionar componente novo.
- Sync ao mudar componente aqui → skill de projeto `boost-claude-sync`
  (`.claude/skills/`).

Instalação (marketplace central `xgodev-plugins`):

```
/plugin marketplace add xgodev/claude-plugin
/plugin install golang-boost@xgodev-plugins
```

A dependência `quality-gate` (gate comparativo pré-push, ver
`docs/development/quality-gate.md`) vive no mesmo marketplace e é resolvida
automaticamente no install.

## Referências (ler quando precisar)

| Doc | Quando |
|---|---|
| `docs/development/gitflow.md` | Issue, branch, commit, PR, fechamento |
| `docs/development/quality-gate.md` | Gate comparativo antes do push |
| `xgodev/boost-claude` → `skills/boost/references/CONTRIBUTING.md` | Adicionar factory/componente |
| `xgodev/boost-claude` → `skills/boost/references/` | Doc de consumo (boot, log, config, factories…) |
| `.claude/skills/boost-claude-sync/SKILL.md` | Sync da doc ao mudar componente |
| README de cada pacote (`factory/`, `wrapper/`, `bootstrap/`, …) | API do componente |
