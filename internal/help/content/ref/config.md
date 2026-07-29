---
title: "ref: config"
summary: "intent.config.yaml: output_dir and named path variables."
---
# ref: config

`intent.config.yaml` is optional; sane defaults apply when it's absent.

```yaml
output_dir: docs                 # where generated markdown is written
paths:
  assets: docs/assets            # {{paths.assets}} resolves to this
  # diagrams: design/diagrams    # your own named vars are welcome
key_case: kebab-case              # or snake_case — the required element-key convention
```

- **`output_dir`** — the root of the generated docs tree (`build` writes here;
  `check` diffs against it). Default `docs`.
- **`paths`** — a map of named path variables referenced in prose as
  `{{paths.<name>}}`. Resolved by **pure literal string substitution** — not a
  template language: no conditionals, no logic, no lookups, no injection surface.
- **`key_case`** — the element-key convention this project requires:
  `kebab-case` (default) or `snake_case`. Enforced by the linter (→ `intent
  help E005`) and the write path (`add`, `record`, `mv --rename`); the schema
  itself accepts either shape. Gen never translates a key's casing — a file or
  directory name is always the key exactly as written, so once a key follows
  this convention its generated path does too.

The generated tree underneath `output_dir`:

```
<output_dir>/product/          jobs / outcomes / …, drs/ (PDRs)
<output_dir>/engineering/      components, drs/ (ADRs)
<output_dir>/change-records/   top-level CRs
```
