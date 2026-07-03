# Design: skill `boost-claude-sync` + hook de lembrete (Iron Law #6)

**Data:** 2026-07-03
**Status:** aprovado por recomendação (usuário AFK na pergunta de enforcement — revisar se discordar)

## Problema

O Iron Law #6 do CLAUDE.md exige que toda mudança de componente neste repo
(`boost`) atualize a skill `boost-*` correspondente em `xgodev/boost-claude`.
A regra existe como texto, mas texto é esquecível: nada garante que a sessão
que edita `factory/contrib/...` lembre de abrir o PR no repo do plugin.

## Decisão

Duas peças, committed neste repo:

1. **Skill de projeto** `.claude/skills/boost-claude-sync/SKILL.md` — o
   workflow completo e determinístico de sincronização (o COMO).
2. **Hook determinístico** (PostToolUse em Edit/Write) — detecta escrita em
   caminho de componente e injeta lembrete apontando para a skill (o
   NÃO-ESQUECER). Registrado em `.claude/settings.json` (compartilhado no
   repo), script em `.claude/hooks/`.

Alternativas descartadas:
- **Só skill**: "garantir" viraria "lembrar"; sem guard-rail.
- **Check no CI**: enforcement só no merge, mais infra (Action + token
  cross-repo); pode ser adicionado depois sem conflitar com esta decisão.

## Escopo do gatilho (o que é "componente")

Caminhos cuja mudança observável exige sync:

- `factory/contrib/<org>/<pkg>/...`
- `wrapper/<x>/...`
- `bootstrap/...`
- `extra/<x>/...`
- `fx/...`
- `model/...`
- `config.go`, `start.go` (raiz — cobertos por `boost-start`)

Mudança puramente interna (refactor sem efeito em API pública, config key,
default ou comportamento documentado) não exige sync — a skill instrui a
avaliar isso explicitamente.

## Workflow da skill (resumo)

1. Identificar componente(s) tocado(s) no diff.
2. Derivar o nome da skill: convenção `boost-<área>-<componente>`
   (ex.: `factory/contrib/labstack/echo` → `boost-factory-echo`;
   `wrapper/log` → `boost-wrapper-log`). Confirmar por busca em
   `skills/` no boost-claude — há irregulares (`boost-start`,
   `boost-fx-modules`, `boost-model-errors`,
   `boost-wrapper-log-backends`, `gqlgen-field-resolvers`).
3. Garantir clone local de `xgodev/boost-claude` em
   `~/Projetos/github.com/xgodev/boost-claude` (clonar se ausente,
   `git pull` se presente).
4. Branch própria no boost-claude; atualizar/criar/remover a skill:
   - componente novo → ler `skills/boost-maintainer/SKILL.md` lá antes;
   - componente removido → remover diretório da skill.
5. Bump de versão em `.claude-plugin/plugin.json` (sem bump o auto-update
   não reconhece).
6. Commit local no boost-claude. **PR não é aberto automaticamente** —
   listado como próximo passo; só executar quando o usuário pedir
   (regra global: no unrequested shared-state actions).

## Hook

- **Evento:** PostToolUse, matcher `Edit|Write`.
- **Script:** `.claude/hooks/boost-claude-sync-reminder.sh` — se o
  `file_path` casar com os caminhos de componente acima, emite
  `additionalContext` lembrando de usar a skill `boost-claude-sync`.
  Não-bloqueante (aviso, não gate); o gate de push continua sendo o
  quality-gate.
- Deduplicação simples por sessão (marker no diretório de sessão/tmp)
  para não spammar a cada edit.

## Teste

Skill segue TDD de `superpowers:writing-skills`: baseline com subagent SEM a
skill (observar o esquecimento), depois com a skill (observar o workflow
correto). Hook testado invocando o script com JSON de exemplo (caminho de
componente vs. caminho neutro).

## Fora de escopo (skills futuras deste repo)

Este é o primeiro item de uma família de skills de padronização do
desenvolvimento no boost (ex.: criar nova factory, gitflow, quality-gate
awareness). Cada uma terá seu próprio ciclo spec → skill.
