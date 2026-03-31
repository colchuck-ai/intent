# Develop component

**`DEV-COMP`**

You are writing a **component specification** under `docs/engineering/components/`—[`component-simple`](../../assets/component-simple.md) or [`component-complex`](../../assets/component-complex.md)—so ownership, interfaces, and behavior are clear enough to build and review against. The engineering tree is chunked **by component**; one pass of this workflow develops **one** component. **In scope:** that component’s spec only. **Out of scope:** other components, or treating this pass as full **architecture** or **component-index** work—those are **Develop architecture** and **Enumerate components**.

**Records:** **ADRs** for **new** architectural decisions at component scope (`…/adrs/`). **CRs** for meaningful edits under that component’s `crs/`. If the engineering root table or overview must change to stay honest, that change follows **`ENG-CR…`** (or a focused edit session) so the root stays in sync.

---

### Draft

Load the engineering root row for this component and the requirements it claims. Write the spec so it **does not silently contradict** the table or product intent. New technical decisions become **ADRs** here when scoped to this component.

### AI-assisted review

Check interface clarity, requirement alignment, and whether sibling components are still consistent. On **traceability**, the row in **`docs/engineering/README.md`** carries the requirement links; the spec should reinforce intent without re-stating that mapping as a second registry. **CRs** for substantive spec edits; **ADRs** for new decisions caught in review.

### Independent review

Human judgment on readiness to build. Same pattern for **CRs** and **ADRs**.
