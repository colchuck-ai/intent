# Enumerate components

**`ENUM-COMP`**

You are adding **component one-liners** and the **components table** in **`docs/engineering/README.md`**: each row ties a component id to a short responsibility and the **product requirements** it helps satisfy (links to requirement paths). **In scope:** the engineering root README only (this is **index work** there). **Out of scope:** creating or editing **`docs/engineering/components/`**—that is **Develop component** (`DEV-COMP`).

You may emit **ADRs** if naming or slicing components implies a **new** architectural decision. **CRs** are normal here—table and overview edits are often substantial (`ENG-CR…`).

---

### Draft

Work from architecture narrative and product requirements so rows are **defensible**. It is acceptable for a requirement to appear on **multiple** rows. Keep the table the **single** place for full requirement ↔ component mapping.

### AI-assisted review

Look for orphan components, missing requirement links, and overlap. Confirm the **components table** remains the authoritative place for full requirement ↔ component mapping—component specs should not quietly replace that with duplicate matrices. Table or narrative changes that matter get **CRs**; new decisions get **ADRs**.

### Independent review

Confirm the decomposition matches how the team builds and owns work. **CRs** and **ADRs** follow the same spirit as earlier sessions.
