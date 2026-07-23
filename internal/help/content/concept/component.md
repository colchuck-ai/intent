---
title: "Component"
summary: "An engineering unit with a single responsibility that fulfills product requirements."
---
# Component

A **component** is an engineering unit of the system with one clear
responsibility. It is where the *how* lives, connected to the *what* by
`fulfills` edges to requirements.

- **Prose field:** `responsibility` (required) — a one-line charter.
- **Lives under:** `engineering.components.<component>`.
- **Required edge:** `fulfills` (≥1) — a product requirement.
- **Optional edge:** `relationships` — a keyed map `<address>: note` to other
  components.
- **Optional structured fields:** `data_model`, `interfaces`, `behavior`,
  `edge_cases`, `success_criteria`.
- **Promotable:** optional `type: inline | document`.

**Why it exists.** Components close the trace graph: every requirement should be
fulfilled by something buildable, and every component should earn its place by
fulfilling a requirement.

**Get it right:** `intent help judgment:altitude` (a component is a charter, not
an outcome and not an implementation dump).
