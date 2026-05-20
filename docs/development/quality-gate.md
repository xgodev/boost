# Quality Gate — boost

boost usa o gate **comparativo compartilhado**
[`github.com/xgodev/quality-gate`](https://github.com/xgodev/quality-gate),
igual ao OpenRig e demais projetos xgodev. Filosofia: o gate quebra quando o
PR **piora** uma métrica vs a base (`develop`). Dívida preexistente nunca
bloqueia — cada PR pode reduzir, nunca aumentar.

> Canônico de uso/contrato: `~/.quality-gate/docs/` após clonar.

## TL;DR

```bash
# primeira vez (clona) — depois, manter atualizado:
git -C ~/.quality-gate pull --ff-only \
  || git clone --depth 1 https://github.com/xgodev/quality-gate.git ~/.quality-gate

~/.quality-gate/qg --base origin/develop
```

Vermelho → arrumar a **causa raiz** do que regrediu → rodar de novo → só
então `git push`.

### Agentes (Claude Code)

Use o **plugin/skill `quality-gate`** — ele clona/atualiza o gate, roda o
dispatcher e interpreta o JSON. Triggers: "rodar quality gate", "rodar QG",
"verificar qualidade", "validar antes do PR", "qa antes do push" (e
equivalentes em EN). É **pré-requisito de qualquer push** (regra 9 do
[`gitflow.md`](gitflow.md)).

## Métricas comparadas (PR vs base) — Go

O gate auto-detecta Go e compara, contra `origin/develop`:

| Métrica | Como conta |
|---|---|
| `fmt` | `gofmt`/`goimports` (rulesets embutidos do gate) |
| `lint` | `go vet` + linter sem warning |
| `build` | `go build ./...` |
| `test` | `go test ./...` (soma de falhas) |
| `complexity` | thresholds default do gate |
| `coverage` | `go test -cover` → `% linhas`, margem `QG_COV_MARGIN` |

> O gate **ignora de propósito** configs de lint/fmt do projeto
> (tamper-resistance) e usa rulesets próprios. Configs locais do repo só
> valem para execução manual local.

## Local vs CI — mesmo dispatcher

| Aspecto | Local | CI |
|---|---|---|
| Comando | `~/.quality-gate/qg --base origin/develop` | mesmo `qg`, clonado no job |
| Baseline | `git archive origin/develop` (cache `/tmp`) | `--baseline-dir` + `--force-full` |
| Falha no CI | — | sticky comment + `request-changes` automático no PR |

## Exit codes

| Código | Significado |
|---|---|
| 0 | Passou / bypass / sem linguagem suportada relevante |
| 1 | Regrediu ≥1 métrica vs base |
| 2 | Erro de ferramenta/setup (NÃO é regressão — relatar stderr) |
| 3 | Nenhuma linguagem suportada detectada |

## Env vars (prefixo `QG_`)

| Variável | Default | Uso |
|---|---|---|
| `QG_BASE_REF` | (vazio) | = `--base`. Vazio → modo absoluto |
| `QG_BASELINE_DIR` | (vazio) | = `--baseline-dir` (CI) |
| `QG_COV_MARGIN` | `1.0` | Tolerância (pp) de coverage |
| `QG_FORMAT` | `text` | `text` ou `json` |
| `QG_BYPASS_REASON` | (vazio) | **NUNCA setar por conta própria.** Força exit 0 + audit log |

## Validação de negócio — testes obrigatórios

Cobertura sozinha não basta. Lógica nova/alterada **exige teste que valide
comportamento esperado**, não só execute o caminho.

- Bug fix → teste vermelho que reproduz o bug primeiro, fix depois (TDD).
- Feature → teste do cenário esperado **e** edge cases.
- Refactor com comportamento observável alterado → teste antes.

## Forbidden

Pra silenciar o gate sem fix real:

- `QG_BYPASS_REASON` por iniciativa própria.
- Editar código/teste/config só pra "passar".
- Marcar testes como `t.Skip` / build tags pra esconder falha.
- `--no-verify` no commit.

A regra: **causa raiz ou escalar.**
