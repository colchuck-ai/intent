# <Architecture Name>

<How components are organized, communicate, and constrain each other.>

## Principles

- <Principle>
- <Principle>

## Constraints

- <Constraint>
- <Constraint>

## Technology Choices

- <Technology>: <Rationale>
- <Technology>: <Rationale>

## Components

<!-- Render each component by its tier (see references/structure.md#child-rendering): -->
<!--   - Tier 0 component — inline block: H3 heading + responsibility line + bold-labeled flat list. -->
<!--   - Tier 1+ component — reference block: H3 heading + responsibility one-liner + "See [link]" line. -->

<!-- Inline form (Tier 0 component) -->
### C<NNN> - <Component Name>

<Responsibility>

**Relationships**

- **C<NNN> - <Component Name>**: <How they communicate or depend>
- **C<NNN> - <Component Name>**: <How they communicate or depend>

<!-- Omit the two record sections below when this component has no ADRs / CRs that name it. -->

**Architectural Decision Records**

- [ADR<NNN> - <ADR Name>](drs/ADR<NNN>-<name>.md)

**Change Records**

- [CR<NNN> - <CR Name>](crs/CR<NNN>-<name>.md)

<!-- Reference form (Tier 1+ component) -->
### C<NNN> - <Component Name>

<Responsibility>

See [C<NNN> - <Component Name>](components/C<NNN>-<name>.md).

## Requirement-Component Map

- **O<NNN>-R<NNN> - <Requirement Name>**: C<NNN> - <Component Name>, C<NNN> - <Component Name>
- **O<NNN>-R<NNN> - <Requirement Name>**: C<NNN> - <Component Name>

## See Also

<!-- Omit any subsection with no entries; omit this entire `## See Also` if every subsection would be empty. Records about the architecture as a whole live in the central directories below. -->

### Architectural Decision Records

- [ADR<NNN> - <ADR Name>](drs/ADR<NNN>-<name>.md)

### Change Records

- [CR<NNN> - <CR Name>](crs/CR<NNN>-<name>.md)
