# Intent

## Elements

### Product

#### Job

#### Outcome

#### Risk

#### Requirement

### Engineering

#### Architecture

## Structure

```txt
docs/
  product/
    crs/
      PROD-CR<NNN>-<name>/
        README.md
      PROD-CR<NNN>-<name>.md
    outcomes/
      O<NNN>-<name>/
        crs/
          O<NNN>-CR<NNN>-<name>/
            README.md
          O<NNN>-CR<NNN>-<name>.md
        pdrs/
          O<NNN>-PDR<NNN>-<name>/
            README.md
          O<NNN>-PDR<NNN>-<name>.md
        requirements/
          O<NNN>-R<NNN>-<name>/
            crs/
              O<NNN>-R<NNN>-CR<NNN>-<name>/
                README.md
              O<NNN>-R<NNN>-CR<NNN>-<name>.md
            pdrs/
              O<NNN>-R<NNN>-PDR<NNN>-<name>/
                README.md
              O<NNN>-R<NNN>-PDR<NNN>-<name>.md
            README.md
          O<NNN>-R<NNN>-<name>.md
        README.md
    pdrs/
      PDR<NNN>-<name>/
        README.md
      PDR<NNN>-<name>.md
    README.md
  engineering/
    adrs/
      ADR<NNN>-<name>/
        README.md
      ADR<NNN>-<name>.md
    components/
      C<NNN>-<name>/
        adrs/
          C<NNN>-ADR<NNN>-<name>/
            README.md
          C<NNN>-ADR<NNN>-<name>.md
        crs/
          C<NNN>-CR<NNN>-<name>/
            README.md
          C<NNN>-CR<NNN>-<name>.md
        README.md
      C<NNN>-<name>.md
    crs/
      ENG-CR<NNN>-<name>/
        README.md
      ENG-CR<NNN>-<name>.md
    README.md
```

## Templates

- product readme
- outcome readme
- requirement readme
- engineering readme
- component readme
- product decision record
- architectural decision record
- change record

## Verbs

- create
- read
- update
- delete

## Phases of a modification
 
1. draft
2. judge
3. review

## Resources

- product
- outcome
- requirement
- architecture
- component

## Review Criteria for Outcomes, Requirements, and Components

- scope is well defined
- overlap is minimized
- when overlapping, elements agree
- their union is complete
- in agreement with their parent (outcomes -> jobs, requirements -> outcome, components -> requirements)
