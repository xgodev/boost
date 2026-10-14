# Registry de erros customizados → código de transporte

## Problema

O mapeamento "erro → código de transporte" está copiado em 3 lugares, cada um
com o mesmo `switch` sobre `errors.Is*`:

- HTTP/Echo — `factory/contrib/labstack/echo/v4/error_handler.go` (`ErrorStatusCode`)
- gRPC — `factory/contrib/google.golang.org/grpc/v1/server/error.go` (`Error`)
- function/CloudEvents — `bootstrap/function/adapter/contrib/cloudevents/sdk-go/v2/core/http/run.go` (`ErrorStatusCode`)

Hoje só dá pra mapear os erros do próprio boost. Não há como dizer que um erro
da aplicação (ex. `XptoError`) deve retornar 404, nem ignorar certos erros. Pra
isso seria preciso editar os 3 switches na mão.

## Objetivo

Um registry único onde a aplicação registra erros customizados e diz qual
**semântico boost** eles têm (que resolve HTTP e gRPC de uma vez), ou marca pra
ignorar. Vale nos 3 transportes.

```go
errors.Register(MyXptoError, errors.KindNotFound) // → 404 HTTP / NotFound gRPC
errors.Ignore(MyNoiseError, errors.IgnoreSilenceLog) // não loga, status normal
```

## Não-objetivos

- Não declarar código HTTP/gRPC cru por erro (sempre via `Kind` semântico).
- Não mudar a semântica dos `Is*` existentes (retrocompat).
- Não tocar `http2`/`graphql-go` (não mapeiam erro hoje).

## Decisões travadas

- Registro por **semântico boost** (`Kind`), não por código cru.
- Escopo: os **3 transportes**.
- Ignore é **configurável por erro**: virar sucesso, só silenciar log, ou ambos.

## Desenho

### 1. `model/errors` — registry + classificação (sem dependência de transporte)

`Kind` é um enum que espelha o catálogo atual de `Is*`. O pacote **não** importa
`net/http` nem `grpc/codes` — `Kind` é só um inteiro semântico.

```go
type Kind int

const (
    KindInternal Kind = iota // default / fallback
    KindNotFound
    KindBadRequest
    KindNotValid
    KindConflict
    KindAlreadyExists
    KindForbidden
    KindUnauthorized
    KindServiceUnavailable
    KindNotImplemented
    KindNotProvisioned
    KindNotSupported
    KindNotAssigned
    KindMethodNotAllowed
    KindTooManyRequests
    KindTimeout
)
```

Registro (chamado no boot, antes de servir):

```go
// casa via errors.Is(err, target) — bom pra sentinelas
func Register(target error, kind Kind)
// casa por predicado/tipo (errors.As) — bom pra XptoError concreto
func RegisterMatch(match func(error) bool, kind Kind)

// ignore configurável
type IgnoreOption int
const (
    IgnoreAsSuccess IgnoreOption = 1 << iota // resposta vira sucesso (200/OK)
    IgnoreSilenceLog                          // não loga (ou nível baixo)
)
func Ignore(target error, opts ...IgnoreOption)            // sem opts = ambos
func IgnoreMatch(match func(error) bool, opts ...IgnoreOption)
```

Resolução:

```go
// precedência: matchers registrados (ordem de registro, 1º vence) → Is* embutido → KindInternal
func Classify(err error) Kind
// política de ignore, se houver
func IgnoreOf(err error) (IgnoreOption, bool)
```

Concorrência: registry protegido por `sync.RWMutex`. Registro no `init`/setup do
boot; leituras (`Classify`/`IgnoreOf`) em request time.

### 2. Tabelas `Kind → código` (na borda de cada transporte)

- HTTP: `restresponse.HTTPStatusFor(kind errors.Kind) int` — casa nova
  compartilhada por **echo** e **function/CE** (ambos já importam `restresponse`;
  `restresponse` passa a importar `model/errors`, sem ciclo). Mata a duplicação
  da tabela HTTP.
- gRPC: tabela `Kind → codes.Code` no pacote `grpc/v1/server`.

### 3. Refactor dos 3 handlers

Cada handler passa a:

1. `kind := errors.Classify(err)` → lookup do código por `Kind`.
2. Antes de montar a resposta: `if pol, ok := errors.IgnoreOf(err); ok` →
   - `IgnoreAsSuccess`: HTTP 200/204, gRPC `codes.OK`/`nil`.
   - `IgnoreSilenceLog`: pula/abaixa o log.

Os `*HTTPError` do echo e `validator.ValidationErrors` continuam tratados como
hoje.

### 4. Skills + docs (Iron Law #6)

No mesmo PR: atualizar `skills/boost-model-errors`, `boost-factory-echo`,
`boost-factory-grpc` + READMEs dos pacotes. Bump da versão do plugin em
`.claude-plugin/plugin.json`.

## Testes (TDD — vermelho primeiro)

- `Register(XptoError, KindNotFound)` → 404 (echo), `codes.NotFound` (gRPC),
  404 (CE).
- `RegisterMatch` por tipo concreto.
- `Ignore` → sucesso; `IgnoreSilenceLog` → status normal sem log; ambos.
- Precedência: erro registrado ganha do `Is*` embutido.
- Compat: todos os `Is*` e mapeamentos default inalterados.
- Concorrência: `Classify` sob leitura concorrente.

## Arquivos

Novos:
- `model/errors/registry.go`, `model/errors/registry_test.go`
- `model/errors/kind.go` (enum)

Alterados:
- `model/restresponse/` — `HTTPStatusFor`
- `factory/contrib/labstack/echo/v4/error_handler.go`
- `factory/contrib/google.golang.org/grpc/v1/server/error.go`
- `bootstrap/function/adapter/contrib/cloudevents/sdk-go/v2/core/http/run.go`
- `skills/boost-model-errors`, `skills/boost-factory-echo`, `skills/boost-factory-grpc`
- `.claude-plugin/plugin.json` (bump)
