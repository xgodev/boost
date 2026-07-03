---
name: boost-claude-sync
description: Use when a boost component was created, updated, or removed — any observable change (public API, config key, default, documented behavior, README) under factory/, wrapper/, bootstrap/, extra/, fx/, model/, or root config.go/start.go — before pushing or claiming the task done.
---

# Sync da doc do plugin (xgodev/boost-claude)

## Overview

A doc de consumo do boost vive em `xgodev/boost-claude` como uma única skill
(`skills/boost/`) com folhas em `skills/boost/references/<grupo>/<nome>.md`.
Mudança observável de componente aqui exige atualizar a reference lá, na
mesma sessão. Iron Law #6 deste repo + Iron Laws #2/#3 do boost-claude.

**A parte local é sempre executável agora**: clone/pull, editar reference,
bump de versão, commit. Só push e PR esperam pedido explícito do usuário.
"Deixo descrito para depois" = red flag do CLAUDE.md aberta.

## Quando NÃO sincronizar

Refactor interno puro — sem efeito em API pública, config key, default ou
comportamento documentado. Diga isso explicitamente no relatório final.

## Workflow

1. **Clone canônico**: `~/Projetos/github.com/xgodev/boost-claude`.
   Ausente → `git clone https://github.com/xgodev/boost-claude.git` nesse
   path. Presente → `git checkout main && git pull origin main`.
2. **Mapear** cada path tocado no boost → reference (tabela abaixo).
3. **Branch** nova no boost-claude a partir de `main` (fluxo segue o
   `docs/development/gitflow.md` do boost).
4. **Editar**:
   - Componente **atualizado** → refletir a mudança na folha (config keys,
     defaults, API, exemplos). Exemplos citam boost por import path Go,
     nunca por caminho de arquivo.
   - Componente **novo** → ler `skills/boost/references/CONTRIBUTING.md`
     antes; criar a folha, adicionar linha no arquivo de grupo (ex.
     `references/factory/messaging.md`) e no índice `skills/boost/SKILL.md`.
   - Componente **removido** → remover a folha e toda entrada de índice
     que aponte para ela.
5. **Bump** de `version` em `.claude-plugin/plugin.json` — sem bump o
   auto-update não reconhece a mudança.
6. **Commit local** no boost-claude (mensagem em inglês, why-focused).
   Sync rotineiro de `references/*.md` NÃO exige o gate
   `superpowers:writing-skills`; autorar/editar um `SKILL.md` exige.
7. **Não** abrir PR nem push por conta própria — listar como próximo passo
   e executar só quando o usuário pedir.

## Mapeamento path → reference

Base: `skills/boost/references/` no boost-claude.

| Path no boost | Reference |
|---|---|
| `factory/contrib/<vendor>/<lib>/vN/` | `factory/<lib>.md` (ex.: `go-redis/v9` → `redis.md`, `labstack/echo` → `echo.md`) |
| `wrapper/<x>/` | `wrapper/<x>.md`; contribs de log (`wrapper/log/contrib/...`) → `wrapper/log-backends.md` |
| `bootstrap/<x>/` | `bootstrap/<x>.md` (adapters → `bootstrap/adapter-<nome>.md`) |
| `extra/<x>/` | `extra/<x>.md` |
| `fx/modules/` | `fx/modules.md` |
| `model/errors/` | `model-errors.md` |
| `config.go`, `start.go` (raiz) | `start.md` |
| plugins de componente (Datadog/OTel/Prometheus) | `plugins.md` + folha do componente |

Irregulares existem (`gcp-api.md`, `net-http2.md`, `gocloud-pubsub.md`…):
na dúvida, `ls`/grep em `skills/boost/references/` pelo nome do componente
— a folha existente ganha da convenção.

## Erros comuns

| Erro | Correção |
|---|---|
| "Sem clone local / rede, sincronizo depois" | Clone é um comando; edit+bump+commit são locais. Fazer agora, nesta sessão. |
| Atualizar a folha e esquecer o bump do `plugin.json` | Bump em TODA mudança de conteúdo (Iron Law #3 do boost-claude). |
| Componente novo sem linha no grupo e no índice `SKILL.md` | Folha órfã nunca é carregada — os 3 pontos são obrigatórios. |
| Citar caminho de arquivo do boost na reference | Só import path Go (`github.com/xgodev/boost/...`). |
| Abrir PR no boost-claude sem o usuário pedir | Commit local + listar PR como próximo passo. |
