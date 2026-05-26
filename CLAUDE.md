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
- `skills/boost-maintainer` — guia de manutenção (criar nova factory/skill).
  Leia antes de adicionar componente.

O plugin `golang-boost` é distribuído pelo marketplace **`xgodev-boost`**
(declarado em `.claude-plugin/marketplace.json` deste repo):

```
/plugin marketplace add xgodev/boost
/plugin install golang-boost@xgodev-boost
```

`marketplace.json` declara `name: "xgodev-boost"` (único globalmente). Há
também o umbrella `xgodev/claude-plugin` (marketplace `xgodev`) que re-lista
este plugin como `boost@xgodev` — usar um caminho ou outro, não os dois.

Dependências declaradas em `.claude-plugin/plugin.json` — instalar
`golang-boost` puxa automaticamente:

- `quality-gate@xgodev` — gate comparativo pré-push (ver
  `docs/development/quality-gate.md`), vive no marketplace `xgodev` (umbrella
  `xgodev/claude-plugin`).

Pré-requisito: o marketplace da dep precisa estar adicionado antes do install
(`/plugin marketplace add xgodev/claude-plugin`); sem ele a dep fica
unresolved e o plugin é desabilitado com `dependency-unsatisfied`.

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
