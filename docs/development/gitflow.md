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
9. **Quality gate verde é pré-requisito de qualquer `git push`.** Ver
   [`quality-gate.md`](quality-gate.md). Vermelho → corrige a causa raiz →
   roda de novo → só então push.
10. **Push após cada commit lógico** (depois do gate verde). Bugfix/hotfix
    mergeia rápido; feature aguarda review. **Nunca mergear feature→develop
    sem o usuário pedir.**

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
