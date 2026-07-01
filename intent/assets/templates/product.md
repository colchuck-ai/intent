# <Product Name>

<What job(s) it serves and for whom.>

## Jobs

### <Job Name>

> When [Situation], I want to [Goal], so I can [Benefit].

<!-- Render each outcome under this job by its tier (see references/structure.md#child-rendering): -->
<!--   - Tier 0 outcome — inline block: H4 heading + minimal pattern + bold-labeled flat lists. -->
<!--   - Tier 1+ outcome — reference block: H4 heading + minimal-pattern one-liner + "See [link]" line. -->

<!-- Inline form (Tier 0 outcome) -->
#### O<NNN> - <Outcome Name>

[Direction] [Metric] [Object] [Context]

**Risks**

- **O<NNN>-RSK<NNN>** - <Risk Name>: [Condition/Event] [Negative Impact on Outcome]
- **O<NNN>-RSK<NNN>** - <Risk Name>: [Condition/Event] [Negative Impact on Outcome]

**Requirements**

- **O<NNN>-R<NNN>** - <Requirement Name>: [Product/Solution] must [Capability/Constraint]
- **O<NNN>-R<NNN>** - <Requirement Name>: [Product/Solution] must [Capability/Constraint]

**Risk-Requirement Map**

- **O<NNN>-RSK<NNN> - <Risk Name>**: O<NNN>-R<NNN> - <Requirement Name>, O<NNN>-R<NNN> - <Requirement Name>
- **O<NNN>-RSK<NNN> - <Risk Name>**: O<NNN>-R<NNN> - <Requirement Name>

<!-- Omit the two record sections below when this outcome has no PDRs / CRs that name it. -->

**Product Decision Records**

- [PDR<NNN> - <PDR Name>](drs/PDR<NNN>-<name>.md)

**Change Records**

- [CR<NNN> - <CR Name>](crs/CR<NNN>-<name>.md)

<!-- Reference form (Tier 1+ outcome) -->
#### O<NNN> - <Outcome Name>

[Direction] [Metric] [Object] [Context]

See [O<NNN> - <Outcome Name>](outcomes/O<NNN>-<name>.md).

## See Also

<!-- Omit any subsection with no entries; omit this entire `## See Also` if every subsection would be empty. Records about the product as a whole live in the central directories below. -->

### Product Decision Records

- [PDR<NNN> - <PDR Name>](drs/PDR<NNN>-<name>.md)

### Change Records

- [CR<NNN> - <CR Name>](crs/CR<NNN>-<name>.md)
