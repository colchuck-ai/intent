---
title: "ref: commands"
summary: "The intent command surface, grouped by purpose."
---
# ref: commands

All commands take `-f, --file <path>` (default `intent.yaml`). Every write runs
`validate` first and re-serializes canonically, so a write is either valid or
refused.

**Explore**

| Command | Does |
|---|---|
| `find [query]` | search name/statement/story/detail; `--type` / `--domain` |
| `tree` | the whole outline |
| `show <addr>` | an element + its outgoing edges ("what is this?") |
| `affects <addr>` | incoming refs ("what breaks if I change this?") |
| `trace <addr>` | the neighborhood; `--depth N` / `--up` / `--down` / `--edges` |

**Write**

| Command | Does |
|---|---|
| `add <type> <parent> <key>` | create an element (`--field …`) |
| `set <addr> <field> <value>` | set one scalar field; `type` cascades |
| `link <kind> <from> <to>` / `unlink` | manage a declared edge |
| `promote <addr>` | make it a document (cascades ancestors) |
| `mv <addr> [new-parent]` | re-parent / `--rename` / `--before` / `--after` |
| `rm <addr>` | delete (refused if it would dangle a reference) |
| `record <pdr\|adr\|cr> <key>` | create a decision/change record |

**Validate / build**

| Command | Does |
|---|---|
| `validate` | run the linter (E001–E004); auto-runs before every write |
| `build [--out docs/]` | generate the markdown docs |
| `check` | build to a temp dir and diff vs on-disk docs; nonzero on drift |

**Help**

| Command | Does |
|---|---|
| `help` | task-first index |
| `help <slug>` | one topic (e.g. `concept:requirement`) |
| `help E0NN` | one error's cause + fix |
| `help --list` / `--json` | flat / machine-readable index |
