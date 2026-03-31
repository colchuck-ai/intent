# Remove component

**`REMOVE-COMP`**

You are **retiring** a component: deleting or archiving its spec and **cleaning** references so the doc set does not lie. **Ripple** includes **`docs/engineering/README.md`** (remove the row, fix links), **sibling components** that referenced it, **ADRs** that might need superseding notes, and **product** tables if requirements were only satisfied by this component—those requirements may need **re-homing** or **rewriting** elsewhere.

**Context span:** start from the component path, then the engineering root, then any **requirement rows** pointing here, then **cross-links** in other components. There is no single recipe; the depth follows the **cascade** of the removal.

**Records:** **CRs** that explain what disappeared and why. **ADRs** if removal **is** an architectural decision (e.g. “we no longer use X”). If you split responsibility across other components, those specs and the table often change in the same **story**—trace with CRs.

---

### Draft

Decide what “gone” means for the system narrative, update the table and specs, and remove or redirect files as your repo conventions require. **Examples:** a merged component absorbs behavior—update the survivor’s spec and CR both sides; a retired service might leave an ADR that records the decision to shut it down. Write **CR** and ADR copy in **[plain language](../guides/plain-language.md)**—short sentences, concrete terms—so the retirement story is easy to follow.

### AI-assisted review

Hunt for orphan links, stale requirement rows, and dangling references in sibling docs. **CRs** and **ADRs** to match the removal story. Where prose is dense or performative, plain it ([`plain-language.md`](../guides/plain-language.md)).

### Independent review

Human check that nothing still claims the old component exists, and that **CR**s and touched specs read clearly for a junior engineer ([`plain-language.md`](../guides/plain-language.md)). **CRs** if the gate finds more cleanup.
