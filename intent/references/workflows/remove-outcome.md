# Remove outcome

**`REMOVE-OUT`**

You are **removing** an **outcome** and everything that **depends** on it: the **outcome folder** (including `requirements/` beneath it), **one-liners** on the **product root**, and **engineering** rows driven by requirements that **only** lived under this outcome. The cascade is **large** by default—load parents first, then children, then the engineering tree.

**Context span:** `docs/product/README.md` (outcome list), the entire **outcome subtree**, **`docs/engineering/README.md`**, and **components** affected by those requirements. **PDRs** might appear if retirement is a **strategic** choice (sunsetting a product line).

**Records:** **CRs** at product root when the outcome list changes; **CRs** on engineering when the table or specs move; **PDRs** when the removal is a **decision** worth recording.

---

### Draft

**Example:** remove the outcome one-liner, then delete or archive the folder, then walk the engineering table for requirement links that no longer exist. **Example:** if requirements migrated elsewhere before deletion, this workflow might be mostly cleanup—still trace with **CRs**. Keep explanations and **CR** narratives in **[plain language](../guides/plain-language.md)**—plain words, short sentences—so the cascade is understandable.

### AI-assisted review

Look for dangling links, empty table rows, and product copy that still mentions the outcome. **CRs** and **PDRs** as the story needs. Tighten any remaining prose that reads clever or academic ([`plain-language.md`](../guides/plain-language.md)).

### Independent review

Human judgment that the product story still makes sense, nothing in engineering still assumes the old outcome, and the updated docs read clearly for a junior engineer ([`plain-language.md`](../guides/plain-language.md)).
