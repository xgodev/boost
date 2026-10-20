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
- `git push` sem quality-gate verde (suspenso por enquanto, ver `docs/development/quality-gate.md`)
- Issue ou PR ausentes (ver `docs/development/gitflow.md`)
- Mudança sem README/site atualizados ou sem issue no `xgodev/claude-plugin`
- Tag, GitHub Release, PR `develop` → `main` ou bump de versão do site/README
  sem o mantenedor mandar lançar. Issue para no merge em `develop`; release
  só por `release/X.Y.Z` (minor/major) ou `bugfix/X.Y.Z` / `hotfix/X.Y.Z`
  (patch), sem `v` no nome da branch (ver `docs/development/gitflow.md`, seção Versão)

## Regras gerais de código

- **`gofmt`/`goimports` limpo. `go vet` e linter sem warning.**
- **`go build ./...` e `go test ./...` verdes** antes de qualquer push.
- **Single source of truth** — config keys e defaults declarados uma vez no
  `config.go` do componente.
- **Separação de concerns** — factory constrói, wrapper abstrai, app consome.
- **Premissa: toda mudança leva a doc junto.** Coisa nova, alterada ou
  removida (componente, config key, anotação, comportamento) → no mesmo PR,
  o README do pacote, o `README.md` da raiz e o site quando eles mostram
  aquilo; e uma issue no `xgodev/claude-plugin` pedindo a atualização da
  skill do boost (ver [Doc do plugin](#doc-do-plugin-fora-deste-repo)).
  Sem as três, a issue não está pronta.
- **`examples/stock` é o app de exemplo completo.** Recurso novo no boost
  (factory, plugin, adapter, anotação do builder) entra nele no mesmo PR,
  com teste em `cmd/<app>/main_test.go` e linha no `examples/stock/README.md`
  (chamada com a resposta real, chave de config, arquivo do recurso).
- **`README.md` da raiz espelha o site** (`docs/site/index.html`): mesmas
  seções, texto, números, exemplo e versão. Mudou um, muda o outro no mesmo PR.
- Teste valida comportamento, não só cobre linha. Bug fix → teste vermelho
  primeiro (TDD).

## Doc do plugin (fora deste repo)

A doc de consumo do boost e o Quality Gate são distribuídos pelo plugin
**`claude-plugin@xgodev`** (repo `xgodev/claude-plugin`). Este repo (`boost`)
**não é marketplace, não contém skills do plugin e não edita o plugin.**
Toda mudança no boost abre uma issue no `xgodev/claude-plugin` (label
`documentation`) dizendo o que mudou, com link para a issue/PR do boost e
para o README que tem os detalhes; quem atualiza a skill é o
`xgodev/claude-plugin`.

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
