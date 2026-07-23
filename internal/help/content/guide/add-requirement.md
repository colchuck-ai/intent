---
title: "guide: add a requirement"
summary: "Name the risk, add the requirement that mitigates it, then wire dependencies."
---
# guide: add a requirement

**Goal:** add a requirement that genuinely secures an outcome.

1. **Check it's a requirement, not a task.** A requirement is a condition that
   must hold — `intent help judgment:requirement-vs-task`.
2. **Make sure the risk exists.** A requirement `mitigates` a risk in the same
   outcome. If the risk isn't there yet, add it first:
   `intent add risk <outcome> <key> --name "…" --statement "…"`
   (`intent help judgment:risk-vs-feature`).
3. **Add the requirement with its mitigates edge** (required, ≥1; repeatable):
   ```
   intent add requirement <outcome> <key> \
     --name "…" --statement "… must …" \
     --mitigates <risk-suffix>
   ```
4. **Wire dependencies** on other requirement leaves, if any:
   `intent link dependsOn <this> <other-requirement>`.
5. **Optional acceptance criteria:** `--acceptance "…"` at add time, repeatable.

The add is validated and re-serialized automatically; a requirement with no
`mitigates` is refused. The command footer points you at the judgment test to
apply while it's fresh.
