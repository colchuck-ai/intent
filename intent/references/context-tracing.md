# Context Tracing

Trace context by **element** (which parts of the model to load together). File paths, templates, and the `docs/` tree are defined in [Files](../SKILL.md#files) and [Structure](structure.md); this section only names elements and relationships.

Use **vertical** tracing along the spine: product, outcome, risk, requirement, architecture, component. Use **horizontal** tracing for lateral context. For **outcome** and **requirement**, **Horizontal** uses **self** (lateral links recorded on that element), **siblings** (peers under the same parent), and **cousins** (same depth, different branch — e.g. elements under a sibling of your parent). For **component**, **Horizontal** uses only **self** and **siblings** — components have no cousins. Risks and risk–requirement links live on the **outcome** element. Requirement–component mapping and system structure live on the **architecture** element. **Do not** load record types (CRs, PDRs, ADRs) until you need that history.

## Product

### Vertical

- **self**
  - The product element.
- **ancestors**
  - (none)
- **descendants**
  1. Each outcome element for this product.

### Horizontal

- Singleton — only one product element; no horizontal trace.

## Outcome

### Vertical

- **self**
  - This outcome element.
- **ancestors**
  1. The product element.
- **descendants**
  1. Each requirement element belonging to this outcome.
  2. The architecture element (to see which components satisfy those requirements).
  3. Each component element that satisfies those requirements.

### Horizontal

- **self**
  1. Risk–requirement mapping and any other lateral links recorded on this outcome (many-to-many among requirements under this outcome).
- **siblings**
  1. Other outcome elements under the same product (compare scope, users, or dependencies).
- **cousins**
  1. Requirement elements under a sibling outcome (when work spans outcomes or shared components).

## Requirement

### Vertical

- **self**
  - This requirement element.
- **ancestors**
  1. The owning outcome element.
  2. The product element.
- **descendants**
  1. The architecture element (for requirement–component mapping).
  2. Each component element that satisfies this requirement.

### Horizontal

- **self**
  1. Dependencies and cross-references recorded on this requirement.
- **siblings**
  1. Other requirement elements under the same outcome (depends, conflicts, shared risk or component — use the parent outcome as needed).
- **cousins**
  1. Requirement elements under a sibling outcome (use the architecture element and relevant outcomes to identify them).

## Architecture

### Vertical

- **self**
  - The architecture element.
- **ancestors**
  1. The product element.
  2. Each outcome and requirement element in scope for your trace (per the requirement–component mapping on the architecture element).
- **descendants**
  1. Each component element the architecture defines.

### Horizontal

- Singleton — only one architecture element; no horizontal trace.

## Component

### Vertical

- **self**
  - This component element.
- **ancestors**
  1. The architecture element.
  2. Each requirement element this component fulfills (per the requirement–component mapping).
  3. Each outcome element that owns those requirements.
  4. The product element.
- **descendants**
  - (none) — components are leaves of this vertical trace.

### Horizontal

- **self**
  1. Relationships, interfaces, and behavior recorded on this component.
- **siblings**
  1. Other component elements that satisfy the same requirement(s) or are named as peers (use the architecture element's mapping and this component's relationships).
