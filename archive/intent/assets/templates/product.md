<!--
  Section optionality — omit empty sections per references/structure.md ("Never emit an empty heading").

  - Title + intro paragraph: required
  - ## Jobs: required (at least one job)
  - ### <Job Name> + job pattern: required per job
  - #### O<NNN> outcome blocks: required per outcome (inline or reference form per child tier — see references/structure.md#child-rendering)
  - **Risks** / **Requirements** / **Risk-Requirement Map** under inline outcomes: omittable when none
  - **Product Decision Records** / **Change Records** under inline outcomes: omittable when none
  - ## See Also: omittable when every subsection would be empty (each subsection omittable individually when empty)
-->
# <Product Name>

<What job(s) it serves and for whom.>

## Jobs

<!-- required (at least one job) -->

### <Job Name>

<!-- required per job -->

> When [Situation], I want to [Goal], so I can [Benefit].

<!-- Render each outcome under this job by its tier (see references/structure.md#child-rendering): -->
<!--   - Tier 0 outcome — inline block: H4 heading + minimal pattern + bold-labeled flat lists. -->
<!--   - Tier 1+ outcome — reference block: H4 heading + minimal-pattern one-liner + "See [link]" line. -->
<!-- Pick exactly one form per outcome — never emit both blocks below for the same child. -->

<!-- Inline form (Tier 0 outcome) — use this block OR the reference block, not both -->
#### O<NNN> - <Outcome Name>

<!-- required per outcome -->

[Direction] [Metric] [Object] [Context]

**Risks**

<!-- omittable when none -->

- **O<NNN>-RSK<NNN>** - <Risk Name>: [Condition/Event] [Negative Impact on Outcome]
- **O<NNN>-RSK<NNN>** - <Risk Name>: [Condition/Event] [Negative Impact on Outcome]

**Requirements**

<!-- omittable when none -->

- **O<NNN>-R<NNN>** - <Requirement Name>: [Product/Solution] must [Capability/Constraint]
- **O<NNN>-R<NNN>** - <Requirement Name>: [Product/Solution] must [Capability/Constraint]

**Risk-Requirement Map**

<!-- omittable when there are no risks, no requirements, or nothing to map -->

- **O<NNN>-RSK<NNN> - <Risk Name>**: O<NNN>-R<NNN> - <Requirement Name>, O<NNN>-R<NNN> - <Requirement Name>
- **O<NNN>-RSK<NNN> - <Risk Name>**: O<NNN>-R<NNN> - <Requirement Name>

**Product Decision Records**

<!-- omittable when none -->

- [PDR<NNN> - <PDR Name>](drs/PDR<NNN>-<name>.md)

**Change Records**

<!-- omittable when none -->

- [CR<NNN> - <CR Name>](crs/CR<NNN>-<name>.md)

<!-- Reference form (Tier 1+ outcome) — use this block OR the inline block above, not both -->
#### O<NNN> - <Outcome Name>

<!-- required per outcome -->

[Direction] [Metric] [Object] [Context]

See [O<NNN> - <Outcome Name>](outcomes/O<NNN>-<name>.md).

## See Also

<!-- omittable when every subsection would be empty; omit each subsection individually when empty. Records about the product as a whole live in the central directories below. -->

### Product Decision Records

- [PDR<NNN> - <PDR Name>](drs/PDR<NNN>-<name>.md)

### Change Records

- [CR<NNN> - <CR Name>](crs/CR<NNN>-<name>.md)
