# Gitflow — boost

```
Issue → Branch (de develop) → Commits → PR → Review/Merge
```

| Branch | Propósito | Merge into |
|---|---|---|
| `main` | Releases | — |
| `develop` | Próxima release | `main` |
| `feature/*` | Funcionalidades | `develop` |
| `bugfix/*` | Correções | `develop` |
| `hotfix/*` | Urgências em produção | `main` + `develop` |
| `release/*` | Preparação de release | `main` + `develop` |

## Regras

1. **Issue sempre, primeiro.** Toda mudança (feature/bug/docs) tem issue no
   GitHub antes de qualquer código. `gh issue list --search` antes de criar
   (evita duplicata). Label apropriada (`documentation`, `bug`, `enhancement`).
2. **PR sempre.** Nada commitado direto em `develop` ou `main`. Toda mudança
   entra por Pull Request — sem exceção.
3. **Nome de branch: `feature/{N}-slug` / `bugfix/{N}-slug`** (`{N}` = número
   da issue) — convenção já em uso no repo. Antes de criar:
   `git fetch && git branch -a | grep {N}-`.
4. **A partir de `develop` atualizado**: `git checkout develop && git pull --ff-only`.
5. **Mergear `develop` antes de trabalho longo**: `git merge origin/develop`.
6. Commits em **inglês**, foco no **"why"**, **sem `Co-Authored-By`**.
7. **NUNCA `Closes #N` / `Fixes #N`** em commits — fechamento é manual e só a
   pedido do usuário.
8. **NUNCA rebase.** Sempre `git merge`, nunca `git pull --rebase`.
9. **Quality gate verde é pré-requisito de qualquer `git push`.** *Suspenso por
   enquanto:* ver o aviso no topo de [`quality-gate.md`](quality-gate.md). Ver
   [`quality-gate.md`](quality-gate.md). Vermelho → corrige a causa raiz →
   roda de novo → só então push.
10. **Push após cada commit lógico** (depois do gate verde). Bugfix/hotfix
    mergeia rápido; feature aguarda review. **Nunca mergear feature→develop
    sem o usuário pedir.**

## Release

**Issue não gera versão.** A issue termina no merge do PR em `develop`: sem
tag, sem GitHub Release, sem PR `develop` → `main`, sem mexer na versão do
site/README e sem branch `feature/N-release-*`. O `develop` acumula as issues
até a release ser lançada.

Release só quando o mantenedor mandar lançar, sempre por `release/vX.Y.Z`.
Release **não tem issue**: o PR `release/vX.Y.Z` → `main` é o registro dela.


1. `git checkout -b release/vX.Y.Z origin/develop`. Versão: patch para
   fix/deps, minor para feature. Mudança de API pequena vai como minor, com
   nota no release; major (`/v2` no import) só por decisão do mantenedor.
2. Na `release/vX.Y.Z`, só preparação da release: a versão fixa em
   `docs/site/index.html` (barra do topo, hero, diff do `go.mod`, timeline) e
   no `README.md` da raiz (diff do `go.mod`, linha de releases). Nada de
   feature nova nessa branch.
3. PR `release/vX.Y.Z` → `main` (`release: vX.Y.Z`), merge com
   `gh pr merge N --merge`.
4. Tag anotada `vX.Y.Z` no merge commit de `main` e `gh release create vX.Y.Z
   --verify-tag --latest --generate-notes --notes-start-tag <tag anterior>`,
   com uma seção `Behavior / API changes` no topo (`--notes`) quando houver.
5. PR `release/vX.Y.Z` → `develop` (traz o bump de versão de volta), merge, e
   apagar a `release/vX.Y.Z`.

O site (`docs/site/`) publica sozinho a cada push em `develop` que mexe em
`docs/site/` (`.github/workflows/site.yml`): a versão nova aparece no site no
passo 5.

## Rastreabilidade — comentários na issue

A issue é o log de auditoria. Comentar em: plano antes de começar; cada push
(hash + arquivos + build/teste); mudança de plano; problema com evidência;
análise técnica; merge; resumo final. Opções A/B/C ao usuário vão na issue
**antes** da pergunta.

## Fechar issue

Só quando o usuário pedir. `gh issue close <N>` — sem auto-close via commit.

## Labels que excluem das release notes

- `duplicate` — escopo idêntico a outra issue (a duplicata é a mais nova).
- `documentation` quando puramente interno (CI/CD, scripts, planejamento) —
  use também `wontfix`/critério do mantenedor para não vazar em notas.
