---
name: intent
description: >-
  Intent tree lifecycle for docs/product and docs/engineering: creation, modification,
  and removal of product and engineering artifacts (outcomes, requirements, components,
  PDRs, ADRs, CRs). Use when authoring, reviewing, or restructuring intent documentation.
compatibility: >-
  Markdown-only skill: no bundled scripts, no network or system packages required.
  Human-authored docs; templates are not validated by CI in this repo.
---

# Intent

This skill supports **managing the lifecycle of an intent tree**—the paired **product** and **engineering** documentation that hangs under `docs/product/` and `docs/engineering/`. The lifecycle is the same in spirit whether you are **creating** artifacts, **changing** them as understanding grows, or **retiring** them when scope shifts. Workflows in this repo spell out how to move through that lifecycle in practice; this page orients you to **what** lives in the tree and **where** to look next.

## What lives in the tree

**Product** and **engineering** split **intent** from **implementation**, without pretending they are independent. Product answers **who**, **what**, and **why**: the problem, the outcomes that matter, and the requirements that must hold. Engineering answers **where** and **how**: how the system is shaped, which **components** exist, and how they meet the product story.

On the **product** side, the tree is anchored at **`docs/product/README.md`**, which holds the product narrative—positioning, boundaries, and (when you have them) outcome **one-liners**. Beneath that, **outcomes** are the big results the work is supposed to produce; **requirements** break those outcomes into testable commitments. **Product decision records (PDRs)** exist to capture **new** strategic or product-level **decisions**—choices about scope, positioning, or policy—not every edit. **Change records (CRs)** at product scopes exist when you need traceability for a **meaningful edit** to the product root, an outcome, or a requirement. Together, those layers are the **product design** in documentation form: intent made explicit and revisable.

On the **engineering** side, **`docs/engineering/README.md`** carries the architectural story and the **components** table that links buildable pieces back to product requirements. **Components** are the tangible units you design, own, and ship. **Architectural decision records (ADRs)** capture **new** technical or structural **decisions**. **Change records** at engineering scopes serve the same traceability role for meaningful edits to the engineering root or a component—distinct from ADRs, which record decisions, not every diff.

## Directory structure and naming

The nested folders, id patterns, kebab names, and **simple** versus **complex** document shapes (single file vs directory with `README.md`) are all specified in one place:

**[`references/guides/structure.md`](references/guides/structure.md)**

Read that before you invent new paths or filenames so the tree stays consistent and diffable.

## Workflows

Activities—developing or enumerating parts of the tree, modifying them, or removing them—are described as **workflows**, each with a stable code and the same three-session rhythm (draft, then AI-assisted review, then independent review). The full index, with links to every workflow file, is here:

**[`references/workflows/README.md`](references/workflows/README.md)**

`ENUM-REQ` indexes on the outcome README only; requirement bodies live under `requirements/` in `DEV-REQ`.

Templates for the concrete markdown shells live under **[`assets/`](assets/README.md)**.

## Guidance

**Ordering.** It is **strongly recommended** to shape **outcomes** and **requirements** on the product side before you invest heavily in **engineering** specifications and **component** decomposition. Architecture and components should answer a clear intent, not invent it. Exceptions exist; the default is product clarity first, then engineering.

**Product root style.** For **`docs/product/README.md`**, choose **one** framing for the life of that tree: either a **simple** narrative ([`assets/product-root-simple.md`](assets/product-root-simple.md)) or a **jobs-to-be-done** outline ([`assets/product-root-jtbd.md`](assets/product-root-jtbd.md)). Default to **simple** unless job-centered language clearly fits your market. **Do not mix** the two styles in one coherent product tree. If you use JTBD, optional background is in [`references/guides/jtbd-primer.md`](references/guides/jtbd-primer.md); the template itself may omit sections that do not help.
