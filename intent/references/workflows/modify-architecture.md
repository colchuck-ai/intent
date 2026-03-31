# Modify architecture

**`MOD-ARCH`**

You are **changing** the high-level story in **`docs/engineering/README.md`**: architecture, principles, or how the system is framed. **Ripple** often reaches **components**—when the narrative moves, impacted specs and tables usually need to move with it.

**Context span:** engineering root first, then **components** whose responsibility or boundaries are touched, and **product requirements** if intent alignment is in question. This is a wide lens by design.

**Records:** **CRs** for meaningful edits (`ENG-CR…` at root; component `crs/` when child docs change). **ADRs** when the change **is** or **encapsulates** a new architectural decision. A CR can **reference** an ADR when both land together.

---

### Draft

Example: renaming a major boundary might mean updating several component READMEs and the table in one conceptual change. Example: a principle shift might force new **ADRs** and then CRs that point to them. Keep the engineering root narrative in **[plain language](../guides/plain-language.md)**—short sentences, common words, teach-first—as you ripple the change.

### AI-assisted review

Look for drift between root, table, and components. Capture **ADRs** where decisions deserve a durable record; **CRs** for traceability on edits. Where the architecture story reads dense or clever, plain it out ([`plain-language.md`](../guides/plain-language.md)).

### Independent review

Second human on whether the story still matches product intent and team reality, and whether the README reads clearly for a junior engineer ([`plain-language.md`](../guides/plain-language.md)). Same pattern for **CRs** and **ADRs**.
