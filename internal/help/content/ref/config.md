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
```

- **`output_dir`** — the root of the generated docs tree (`build` writes here;
  `check` diffs against it). Default `docs`.
- **`paths`** — a map of named path variables referenced in prose as
  `{{paths.<name>}}`. Resolved by **pure literal string substitution** — not a
  template language: no conditionals, no logic, no lookups, no injection surface.

The generated tree underneath `output_dir`:

```
<output_dir>/product/          jobs / outcomes / …, drs/ (PDRs)
<output_dir>/engineering/      components, drs/ (ADRs)
<output_dir>/change-records/   top-level CRs
```
