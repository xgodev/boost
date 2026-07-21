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

### Red flags — PARAR e reportar

- SDK upstream instanciado direto (sem passar pela factory)
- `os.Getenv` / logger de terceiro fora do wrapper
- API pública alterada sem teste e sem nota de compat
- `git push` sem quality-gate verde
- Issue ou PR ausentes (ver `docs/development/gitflow.md`)

## Regras gerais de código

- **`gofmt`/`goimports` limpo. `go vet` e linter sem warning.**
- **`go build ./...` e `go test ./...` verdes** antes de qualquer push.
- **Single source of truth** — config keys e defaults declarados uma vez no
  `config.go` do componente.
- **Separação de concerns** — factory constrói, wrapper abstrai, app consome.
- Documentação é parte da tarefa: mudou componente/config key/comportamento →
  atualizar o README do pacote no mesmo commit.
- Teste valida comportamento, não só cobre linha. Bug fix → teste vermelho
  primeiro (TDD).

## Doc do plugin (fora deste repo)

A doc de consumo do boost e o Quality Gate são distribuídos pelo plugin
**`claude-plugin@xgodev`** (repo `xgodev/claude-plugin`). Este repo (`boost`)
**não é marketplace, não contém skills do plugin e não sincroniza doc de
plugin** — isso é responsabilidade do `xgodev/claude-plugin`, fora do escopo
daqui.

Instalação do plugin:

```
/plugin marketplace add xgodev/claude-plugin
/plugin install claude-plugin@xgodev
```

## Referências (ler quando precisar)

| Doc | Quando |
|---|---|
| `docs/development/gitflow.md` | Issue, branch, commit, PR, fechamento |
| `docs/development/quality-gate.md` | Gate comparativo antes do push |
| README de cada pacote (`factory/`, `wrapper/`, `bootstrap/`, …) | API do componente |
