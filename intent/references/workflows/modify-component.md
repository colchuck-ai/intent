# Modify component

**`MOD-COMP`**

You are **changing** an existing component spec—behavior, boundaries, interfaces—or adjusting how it is described. The **ripple** might include **`docs/engineering/README.md`** (table, principles, overview) if the change redefines responsibility or requirement links.

**Context span:** load the component artifact, the engineering root, and any **sibling components** or **requirements** that might be implied by the edit. Widen until the change makes sense in isolation and next to its neighbors.

**Records:** **CRs** should track **meaningful** edits—often a **component-scoped** CR and sometimes an **`ENG-CR…`** if the root moves. **ADRs** when you introduce a **new** architectural decision (not for every tweak).

---

### Draft

Make the edit with propagation in mind: if requirement links or the table row change, update the engineering root in the same story. Prefer one coherent narrative over a scattered fix, and keep the spec in **[plain language](../guides/plain-language.md)**—plain words, short sentences—so the ripple is easy to follow.

### AI-assisted review

Check consistency across component, table, and product requirements. The **engineering README** row should still be where full requirement links live; the spec should not sprout a parallel mapping block. Review for plain readable prose ([`plain-language.md`](../guides/plain-language.md)). Emit **CRs** for edits not yet traced; **ADRs** for new decisions.

### Independent review

Human pass for ownership and blast radius, and whether the spec reads clearly for a junior engineer ([`plain-language.md`](../guides/plain-language.md)). **CRs** for meaningful late edits; **ADRs** if decisions surface at the gate.
