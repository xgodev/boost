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
6. **O plugin `golang-boost` co-evolui com o código.** Mudou um componente →
   atualize a skill correspondente em `skills/boost-*` no mesmo PR. Skill nova
   para componente novo. Versão do plugin sobe (`.claude-plugin/plugin.json`)
   quando o conteúdo muda — sem bump, auto-update não reconhece.

### Red flags — PARAR e reportar

- SDK upstream instanciado direto (sem passar pela factory)
- `os.Getenv` / logger de terceiro fora do wrapper
- API pública alterada sem teste e sem nota de compat
- Componente tocado sem a skill `boost-*` correspondente atualizada
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

## Skills do plugin (este repo É o plugin)

- `skills/boost-*` — skills de **consumo** (como usar boost). Distribuídas no
  plugin `golang-boost`.
- `.claude/skills/*` — skills **internas** do projeto (manutenção). Não fazem
  parte do plugin.
- `skills/boost-maintainer` — guia de manutenção (criar nova factory/skill).
  Leia antes de adicionar componente.

## Referências (ler quando precisar)

| Doc | Quando |
|---|---|
| `docs/development/gitflow.md` | Issue, branch, commit, PR, fechamento |
| `docs/development/quality-gate.md` | Gate comparativo antes do push |
| `skills/boost-maintainer/SKILL.md` | Adicionar factory/componente/skill |
| `skills/boost-start/SKILL.md` | Sequência de boot |
| `skills/boost-wrapper-log/SKILL.md` | Logging via wrapper |
| `skills/boost-wrapper-config/SKILL.md` | Config namespacing |
| README de cada pacote (`factory/`, `wrapper/`, `bootstrap/`, …) | API do componente |
