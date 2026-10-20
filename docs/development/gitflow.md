# Gitflow — boost

```
Issue → Branch (de develop) → Commits → PR → merge em develop
                 ... develop acumula ...
Versão → release|bugfix/X.Y.Z → PR → main → tag vX.Y.Z → PR → develop
```

| Branch | Propósito | Sai de | Merge into |
|---|---|---|---|
| `main` | O que foi lançado; cada merge é uma tag | — | — |
| `develop` | Próxima versão | — | só pela branch de versão |
| `feature/{N}-slug` | Issue de funcionalidade | `develop` | `develop` |
| `bugfix/{N}-slug` | Issue de correção | `develop` | `develop` |
| `release/X.Y.Z` | Versão minor/major | `develop` | `main` + `develop` |
| `bugfix/X.Y.Z` | Versão patch | `develop` | `main` + `develop` |
| `hotfix/X.Y.Z` | Versão patch urgente, de produção | `main` | `main` + `develop` |

**Só a tag leva `v`** (`v1.6.2`), porque o Go exige. Branch, milestone,
título de PR e de GitHub Release usam `X.Y.Z` (`1.6.2`). `{N}` é o número da issue; `X.Y.Z` é a versão.

## Regras

1. **Issue sempre, primeiro.** Toda mudança (feature/bug/docs) tem issue no
   GitHub antes de qualquer código. `gh issue list --search` antes de criar
   (evita duplicata). Label apropriada (`documentation`, `bug`, `enhancement`).
   Responsável `jpfaria` (`--assignee jpfaria`) e o milestone aberto da
   próxima versão (`--milestone X.Y.Z`); sem milestone aberto, criar o do
   próximo patch.
2. **PR sempre.** Nada commitado direto em `develop` ou `main`. Toda mudança
   entra por Pull Request — sem exceção.
3. **Uma branch por issue**, `feature/{N}-slug` ou `bugfix/{N}-slug`, a partir
   de `develop` atualizado (`git fetch && git branch -a | grep {N}-` antes de
   criar). Nunca a partir de outra branch de issue: nada de PR empilhado.
4. **Mergear `develop` antes de trabalho longo**: `git merge origin/develop`.
5. Commits e PRs em **inglês**, foco no **"why"**, **sem `Co-Authored-By`**
   nem outro trailer de atribuição.
6. **Nunca `Closes #N` / `Fixes #N`** em commit ou PR: a issue fecha à mão
   (ver [Fechar issue](#fechar-issue)).
7. **Nunca rebase.** Sempre `git merge`, nunca `git pull --rebase`.
8. **Push após cada commit**, com `go build ./...` e `go test ./...` verdes.
   O quality gate antes do push está suspenso (ver
   [`quality-gate.md`](quality-gate.md)); religado, vermelho → corrige a causa
   raiz → roda de novo → só então push.
9. **Issue pronta = PR mergeado em `develop`** (`gh pr merge N --merge`), sem
   esperar review. Depois do merge, apagar a branch remota. O PR já leva o
   README do pacote, o `README.md` da raiz e o site atualizados, e a issue
   tem o link da issue aberta no `xgodev/claude-plugin` pedindo a
   atualização da skill do boost (o boost não edita o plugin).

## Versão (release)

**Issue não gera versão.** A issue termina no merge do PR em `develop`: sem
tag, sem GitHub Release, sem PR `develop` → `main`, sem mexer na versão do
site/README. O `develop` acumula as issues até a versão ser lançada.

Versão só quando o mantenedor mandar lançar, e com o número confirmado por
ele. O número segue o semver e escolhe a branch:

- major (`X` anda, `/vX` no import) ou minor (`Y` anda): `release/X.Y.Z`;
- patch (só `Z` anda): `bugfix/X.Y.Z` saindo de `develop`, ou `hotfix/X.Y.Z`
  saindo de `main` quando é urgência em produção.

Abaixo, `<branch>` é essa branch. A versão **não tem issue**: o PR `<branch>`
→ `main` é o registro dela.

0. Milestone `X.Y.Z` (`gh api repos/xgodev/boost/milestones -f title=X.Y.Z`)
   com todas as issues mergeadas em `develop` desde a tag anterior
   (`gh issue edit N --milestone X.Y.Z`).
1. `git checkout -b <branch> origin/develop` (`origin/main` no hotfix).
2. Na `<branch>`, só a preparação: a versão fixa em `docs/site/index.html`
   (barra do topo, hero, diff do `go.mod`, timeline) e no `README.md` da raiz
   (diff do `go.mod`, linha de releases). Nada de código novo nessa branch;
   no hotfix, a correção entra por uma issue com PR para a `<branch>`.
3. PR `<branch>` → `main` (`release: X.Y.Z`), merge com `gh pr merge N --merge`.
4. Tag anotada `vX.Y.Z` no merge commit de `main` e `gh release create vX.Y.Z
   --title X.Y.Z --verify-tag --latest --generate-notes --notes-start-tag <tag anterior>`,
   com uma seção `Behavior / API changes` no topo (`--notes`) quando houver.
   Fechar o milestone.
5. PR `<branch>` → `develop` (`release: X.Y.Z back to develop`, traz o bump de versão de volta), merge, e
   apagar a `<branch>`.

O site (`docs/site/`) publica sozinho a cada push em `develop` que mexe em
`docs/site/` (`.github/workflows/site.yml`): a versão nova aparece no site no
passo 5.

## Rastreabilidade — comentários na issue

A issue é o log de auditoria. Comentar em: plano antes de começar; cada push
(hash + arquivos + build/teste); mudança de plano; problema com evidência;
análise técnica; merge; resumo final. Opções A/B/C ao usuário vão na issue
**antes** da pergunta.

## Fechar issue

Depois do merge em `develop`: comentar o PR e o que entrou, e
`gh issue close <N>`. Nunca por auto-close de commit ou PR.

## Labels que excluem das release notes

- `duplicate` — escopo idêntico a outra issue (a duplicata é a mais nova).
- `documentation` quando puramente interno (CI/CD, scripts, planejamento);
  `wontfix` para o que não deve aparecer nas notas.
