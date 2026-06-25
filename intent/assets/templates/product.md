# <Product Name>

<What job(s) it serves and for whom.>

## Jobs

### <Job Name>

> When [Situation], I want to [Goal], so I can [Outcome].

#### Outcomes

<!-- Render each outcome by its tier (see references/structure.md#child-rendering): -->
<!--   - Tier 0 outcome → inline form: outcome bullet with Risks / Requirements / Risk-Requirement Map sub-bullets. -->
<!--   - Tier 1+ outcome → reference form: a single link bullet. -->

<!-- Inline form (Tier 0 outcome) -->
- **O<NNN>** - <Outcome Name>: [Verb] [Unit of Measure] [Object]
  - Risks:
    - **O<NNN>-RSK<NNN>** - <Risk Name>: [Condition/Event] [Negative Impact on Outcome]
    - **O<NNN>-RSK<NNN>** - <Risk Name>: [Condition/Event] [Negative Impact on Outcome]
  - Requirements:
    - **O<NNN>-R<NNN>** - <Requirement Name>: [Product/Solution] must [Capability/Constraint]
    - **O<NNN>-R<NNN>** - <Requirement Name>: [Product/Solution] must [Capability/Constraint]
  - Risk-Requirement Map:
    - **O<NNN>-RSK<NNN> - <Risk Name>**: O<NNN>-R<NNN> - <Requirement Name>, O<NNN>-R<NNN> - <Requirement Name>
    - **O<NNN>-RSK<NNN> - <Risk Name>**: O<NNN>-R<NNN> - <Requirement Name>

<!-- Reference form (Tier 1+ outcome) -->
- [O<NNN> - <Outcome Name>](outcomes/O<NNN>-<name>.md)

## Appendix: Decision Records

<!-- Omit when empty. Each entry follows the pdr-inline.md template. -->

## Appendix: Change Records

<!-- Omit when empty. Each entry follows the cr-inline.md template. -->

## See Also

### Product Decision Records

- [<PDR Name>](pdrs/<PDR filename>)

### Change Records

- [<CR Name>](crs/<CR filename>)
