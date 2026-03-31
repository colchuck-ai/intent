# Workflows

Each file is **one activity**: what you are trying to do, told as a story with **three sessions**—**draft**, then **AI-assisted review**, then **independent review**. You can run them back-to-back or spread them across days; the names describe **roles**, not calendar time. In each workflow, **[plain language](../guides/plain-language.md)** is **woven into** those session descriptions—how you draft, what the AI-assisted pass checks, and what the human gate expects—so wording stays junior-friendly alongside traceability and records, not as a standalone block.

In **draft**, you gather enough context to write honestly—usually the artifact in focus and its parents, and sometimes siblings when overlap matters. In **modify** and **remove** workflows, that lens often has to widen to the **ripple** of the change (children, the engineering table, or the whole product tree).

Stable **`WF-*`** codes are for shorthand in conversation (“we’re in **ENUM-REQ** draft”); they are not meant to feel like procedure numbers.

If you need a **rough lifecycle read** (greenfield vs brownfield) using only **`ls` and `grep`**—to bias which workflows are likely—see **[`../guides/intent-tree-phase.md`](../guides/intent-tree-phase.md)**.

**Greenfield (product then engineering, recommended order)**

| Code | Workflow | File |
|------|----------|------|
| `DEV-PROD` | Develop product | [`develop-product.md`](develop-product.md) |
| `ENUM-OUT` | Enumerate outcomes | [`enumerate-outcomes.md`](enumerate-outcomes.md) |
| `DEV-OUT` | Develop outcome | [`develop-outcome.md`](develop-outcome.md) |
| `ENUM-REQ` | Enumerate requirements | [`enumerate-requirements.md`](enumerate-requirements.md) |
| `DEV-REQ` | Develop requirement | [`develop-requirement.md`](develop-requirement.md) |
| `DEV-ARCH` | Develop architecture | [`develop-architecture.md`](develop-architecture.md) |
| `ENUM-COMP` | Enumerate components | [`enumerate-components.md`](enumerate-components.md) |
| `DEV-COMP` | Develop component | [`develop-component.md`](develop-component.md) |

**Change**

| Code | Workflow | File |
|------|----------|------|
| `MOD-COMP` | Modify component | [`modify-component.md`](modify-component.md) |
| `MOD-ARCH` | Modify architecture | [`modify-architecture.md`](modify-architecture.md) |
| `MOD-REQ` | Modify requirement | [`modify-requirement.md`](modify-requirement.md) |
| `MOD-OUT` | Modify outcome | [`modify-outcome.md`](modify-outcome.md) |
| `MOD-PROD` | Modify product | [`modify-product.md`](modify-product.md) |

**Removal**

| Code | Workflow | File |
|------|----------|------|
| `REMOVE-COMP` | Remove component | [`remove-component.md`](remove-component.md) |
| `REMOVE-REQ` | Remove requirement | [`remove-requirement.md`](remove-requirement.md) |
| `REMOVE-OUT` | Remove outcome | [`remove-outcome.md`](remove-outcome.md) |
