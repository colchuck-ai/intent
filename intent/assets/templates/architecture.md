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
<!--   - Tier 0 component — inline block: H3 heading + bold-labeled flat list. -->
<!--   - Tier 1+ component — reference block: H3 heading + single "See [link]" line. -->

<!-- Inline form (Tier 0 component) -->
### C<NNN> - <Component Name>

<Responsibility>

**Relationships**

- **C<NNN> - <Component Name>**: <How they communicate or depend>
- **C<NNN> - <Component Name>**: <How they communicate or depend>

<!-- Reference form (Tier 1+ component) -->
### C<NNN> - <Component Name>

See [C<NNN> - <Component Name>](components/C<NNN>-<name>.md).

## Requirement-Component Map

- **O<NNN>-R<NNN> - <Requirement Name>**: C<NNN> - <Component Name>, C<NNN> - <Component Name>
- **O<NNN>-R<NNN> - <Requirement Name>**: C<NNN> - <Component Name>

## Appendix: Decision Records

<!-- Omit when empty. Each entry follows the adr-inline.md template. -->

## Appendix: Change Records

<!-- Omit when empty. Each entry follows the cr-inline.md template. -->

## See Also

### Architectural Decision Records

- [<ADR Name>](adrs/<ADR filename>)

### Change Records

- [<CR Name>](crs/<CR filename>)
