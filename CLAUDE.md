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
6. **O plugin `golang-boost` vive em `xgodev/boost-claude` e co-evolui.** As
   skills NÃO ficam mais neste repo. Mudou um componente → atualize a skill
   `boost-*` correspondente **em `xgodev/boost-claude`**, em PR próprio lá.
   Skill nova para componente novo. Bump do `plugin.json` daquele repo quando o
   conteúdo muda — sem bump, auto-update não reconhece. (O acoplamento "mesmo
   PR" virou disciplina entre dois repos; sem o PR no boost-claude, a doc fica
   defasada.)

### Red flags — PARAR e reportar

- SDK upstream instanciado direto (sem passar pela factory)
- `os.Getenv` / logger de terceiro fora do wrapper
- API pública alterada sem teste e sem nota de compat
- Componente tocado sem a skill `boost-*` correspondente atualizada em `xgodev/boost-claude`
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

## Skills do plugin (vivem em `xgodev/boost-claude`)

As skills `boost-*` foram extraídas para o repo
**[`xgodev/boost-claude`](https://github.com/xgodev/boost-claude)** — assim
instalar o plugin não clona mais o framework inteiro. Este repo (`boost`) **não
é mais um marketplace**.

- `boost-*` — skills de **consumo** (como usar boost), distribuídas no plugin
  `golang-boost`.
- `boost-maintainer` — guia de manutenção (criar nova factory/skill). Leia em
  `boost-claude` antes de adicionar componente.

Instalação (a partir do novo repo):

```
/plugin marketplace add xgodev/boost-claude
/plugin install golang-boost@xgodev-boost
```

O marketplace continua se chamando `xgodev-boost` e o plugin `golang-boost` —
só mudou a URL do `marketplace add`. O umbrella `xgodev/claude-plugin`
(marketplace `xgodev`) re-lista este plugin como `boost@xgodev`; ao mover, o
`source` lá precisa apontar para `xgodev/boost-claude` (ajuste em outro repo).

Dependência (declarada no `plugin.json` de `boost-claude`):
`quality-gate@xgodev-quality-gate` — gate comparativo pré-push (ver
`docs/development/quality-gate.md`). Pré-requisito antes do install:
`/plugin marketplace add xgodev/quality-gate`; sem ele a dep fica unresolved e
o plugin é desabilitado com `dependency-unsatisfied`.

## Referências (ler quando precisar)

| Doc | Quando |
|---|---|
| `docs/development/gitflow.md` | Issue, branch, commit, PR, fechamento |
| `docs/development/quality-gate.md` | Gate comparativo antes do push |
| `xgodev/boost-claude` → `skills/boost-maintainer/SKILL.md` | Adicionar factory/componente/skill |
| `xgodev/boost-claude` → `skills/boost-*` | Skills de consumo (boot, log, config, factories…) |
| README de cada pacote (`factory/`, `wrapper/`, `bootstrap/`, …) | API do componente |
