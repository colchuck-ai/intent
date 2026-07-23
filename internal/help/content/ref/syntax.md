---
title: "ref: syntax"
summary: "The {{ }} interpolation forms and the three accessors."
---
# ref: syntax

**One interpolation syntax, `{{ }}`,** usable in any prose field, disambiguated
by prefix:

| Prefix | Kind | Addresses |
|---|---|---|
| `paths.*` | config path variable | a directory/file (see `ref:config`) |
| `product.` / `engineering.` / `change_records.` | absolute element address | an element |
| `.` (leading dot) | relative address | an element in the current container |

**Accessors** — the terminal token, explicit (no positional inference):

| Accessor | Yields |
|---|---|
| `.name` | the target's `name` (text) |
| `.path` | the computed destination path (raw string) |
| `.link` | a rendered markdown link to the target, text = its `name` |

`.link` is what most prose uses — it re-derives the correct relative href from
the tree, so links survive an `mv`. Use `.name` / `.path` for hand-composed
markdown, e.g. custom link text: `[the validator]({{.validator.path}})`.

**Assets** are embedded with plain markdown + a path variable — there is no image
accessor: `![flow]({{paths.assets}}/linter-flow.png)`.

Prose links are cosmetic — they create no trace edges — but must still resolve
(a dangling prose address is `E001`, same as a dangling declared edge).
