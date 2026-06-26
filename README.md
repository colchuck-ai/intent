# intent

An [Agent Skill](https://agentskills.io/home) that teaches an AI assistant to write product and engineering docs as one **traceable tree** — so every component traces back to a requirement, every requirement to a risk, and every risk to a real customer outcome. No more docs that drift from why the work exists.

## What you get

Two spines that connect product intent to how it's built:

- **Product** — Product, Jobs, Outcomes, Risks, Requirements
- **Engineering** — Architecture, Components

Plus **records** that keep the story coherent over time:

- **CR** (Change Record) — what changed and why
- **PDR** (Product Decision Record) — a product decision, its options, and consequences
- **ADR** (Architectural Decision Record) — the same, for engineering

Everything carries a stable ID (e.g. `O001-R002`) so links survive even when there's no file yet.

## Starts small, grows only when needed

You don't scaffold a directory tree up front. Every element starts as a one-line statement on its parent doc and is promoted only when a concrete need appears:

| Tier | On disk | Promote when… |
|---|---|---|
| **0 — inline** | a line on the parent document | (default) |
| **1 — simple** | its own `<slug>.md` file | it needs detail, edge cases, or its own records |
| **2 — nested** | a `<slug>/` folder | it needs child files or folders |

A minimal product is just two files (`docs/product/README.md`, `docs/engineering/README.md`); a complex one nests as far as it needs to. Records follow the same ladder (none, inline, simple, nested).

## How to use it

This is an Agent Skill — point your assistant at it and ask it to create, trace, or review intent docs. It handles IDs, paths, tiers, and a draft, judge, review workflow for changes.

Install it (defaults to Cursor, `~/.cursor/skills/intent/`):

```bash
curl -fsSL https://raw.githubusercontent.com/colchuck-ai/intent/main/install.sh | bash
```

For another agent, pass an `--agent` (`cursor`, `claude`, `agents`) or a target directory:

```bash
# Claude install location: ~/.claude/skills/intent
curl -fsSL https://raw.githubusercontent.com/colchuck-ai/intent/main/install.sh | bash -s -- --agent claude

# Custom directory
curl -fsSL https://raw.githubusercontent.com/colchuck-ai/intent/main/install.sh | bash -s -- ~/path/to/skills/intent
```

Or copy the [`intent/`](intent/) directory into your agent's skills directory manually.

The full rules, element definitions, paths, and templates live in [`intent/SKILL.md`](intent/SKILL.md).
