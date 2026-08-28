# Command Resurface — Requirements

The target command surface: `verb thing`. Replaces the flat verb surface in
`internal/cli/root.go:37`.

Status: settled. All four open questions decided (§9). Sequencing lives in
[implementation-plan.md](implementation-plan.md).

## 1. The grammar

```
intent <verb> [<thing>] [<positional>] [flags]
```

One syntactic slot for the thing. The noun/subnoun distinction survives as a **taxonomy in the
help**, not as nesting in the parser (§3). Three rules govern everything below.

**R1 — The thing is a type-scope, not decoration.**
`Resolve` (`internal/tree/address.go:37`) currently declares `AmbiguousError` whenever a
whole-segment suffix matches more than one element, regardless of kind. Once the verb carries a
thing, the thing filters candidates *before* ambiguity is judged. `intent show requirement auth`
resolves even when a component is also keyed `auth`. This is the payoff of the redesign: bare keys
become usable, because the grammar supplies the disambiguation the address doesn't.

**R2 — Positionals name the subject. Flags name everything the subject relates to.**
The subject is a new key, an existing locator, or a line of prose. Parents, endpoints, siblings, and
field values are all flags. Every command therefore has at most one bare positional — which is what
keeps the surface scriptable.

**R3 — A locator flag is named after the kind it locates; `--to` is the fallback.**
`--job`, `--outcome`, `--component`. The flag name is what scopes the locator's resolution, so the
name is load-bearing. Where a kind name can't be used — because both endpoints of a relationship are
the same kind, or the target's kind is open — the flag is `--to`, which marks direction explicitly at
exactly the point where direction is what's at stake.

### Why R2/R3 are cheap

Every child kind has exactly one legal parent kind (`internal/tree/tree.go:143-265`), and all but
three are parented to a singleton root:

| element | parent | locator flag needed? |
|---|---|---|
| job | product | no — singleton |
| outcome | job | **`--job`** |
| risk | outcome | **`--outcome`** |
| requirement | outcome | **`--outcome`** |
| component | architecture | no — singleton |
| principle | architecture | no — singleton |
| constraint | architecture | no — singleton |
| pdr | product | no — singleton |
| adr | architecture | no — singleton |
| cr | root | no — singleton |

Three parent flags across the whole surface. The grammar is not ceremonious.

The alignment with re-parenting is exact: only outcomes, risks, and requirements can be re-parented
(`internal/cli/mv.go:18-19`) — precisely the three that have a parent flag. So `--job` means "put it
here" in both `add` and `set`.

## 2. Verbs

| verb | shape | replaces |
|---|---|---|
| `add` | `add <thing> [<positional>] [flags]` | `add`, `record`, `link` |
| `rm` | `rm <thing> <loc-or-selector>` | `rm`, `unlink` |
| `set` | `set <element> <loc> [flags]` | `set`, `mv` |
| `show` | `show <element> <loc>` | `show` |
| `list` | `list [<thing>] [query] [flags]` | `find` |
| `tree` | `tree [<element> <loc>]` | `tree` |
| `trace` | `trace <element> <loc> [flags]` | `trace` |
| `impact` | `impact <element> <loc>` | `affects` — **redefined**, see below |
| `promote` | `promote <element> <loc>` | `promote`, `set … type document` |
| `demote` | `demote <element> <loc> [--cascade]` | `set … type inline` |
| `validate` | `validate` | `validate` |
| `build` | `build [--check] [--out]` | `build`, `check` |
| `install` | `install skill --agent <a> [--dir <d>]` | `install-skill` |
| `help` | `help [slug] [--list] [--json]` | `help` |

`validate` and `build` take no thing: they are whole-tree operations, as are bare `list` and bare
`tree`. That is a coherent second class, not an exception.

### `impact` is the transitive closure — not a rename of `affects`

`affects` today is `ix.Incoming(addr)` (`internal/cli/affects.go:27`), which is exactly
`Trace(addr, 1, up, !down, {})` (`internal/tree/query.go:76`) — direct refs only. That is the wrong
default for impact analysis: a change's real reach is transitive, and this repo has a product job
named *make a surgical change with confidence*.

So `impact` walks the **full transitive incoming closure**, and that is what distinguishes it from
`trace`:

```
intent impact requirement token-ttl               # everything that transitively references it
intent trace  requirement token-ttl --up --depth 2  # bounded neighborhood walk
```

Implementation: `Trace` takes `maxDepth int` and clamps negatives to 0
(`internal/tree/query.go:74-76`). `impact` needs an unbounded mode — either a sentinel
(`maxDepth < 0` meaning unbounded, which changes the current clamp) or a dedicated
`Index.Impact(addr)`. Prefer the dedicated method: the clamp is load-bearing for `trace` and
shouldn't grow a second meaning.

### `build --check`, `promote` / `demote`

`check` folds into `build --check`, matching `gofmt -l` / `prettier --check`: same `--out` resolution
(`internal/cli/build.go:65`), same drift report, same non-zero exit, no write.

`promote` and `demote` own the `inline | document` axis outright. `set` loses its `type` field and its
`--cascade` flag; `--cascade` moves to `demote`, keeping the refusal-unless-cascade rule
(`internal/cli/set.go:16-21`). `add` gains a `--document` boolean on the three promotable kinds.

### `trace --edges` becomes `trace --via`

`--edges` (`internal/cli/trace.go:50`) takes edge-kind names. Those names are now the relationship
words of §6, and "edge" has left the user-facing vocabulary — so the flag reads
`trace requirement token-ttl --via mitigation,fulfillment`.

## 3. The `add` namespace: one slot, three classes

`rel` is not a noun. Every thing is a single word after the verb, and the taxonomy lives in the help.

| class | members | addressable? |
|---|---|---|
| **Elements** | `product` `job` `outcome` `risk` `requirement` `architecture` `component` `principle` `constraint` `pdr` `adr` `cr` | yes — have keys and dotted addresses |
| **Parts** | `criterion` `edge-case` `option` `consequence` | no — belong to one element |
| **Relationships** | `mitigation` `dependency` `fulfillment` `collaboration` `effect` | no — connect two elements |
| **Tool** | `skill` | n/a — only `install` |

Nineteen words share the `add` namespace (10 addable elements + 4 parts + 5 relationships). Cobra
v1.10.2 has native `AddGroup`/`GroupID`, so `intent add --help` prints them grouped:

```
Elements:
  job  outcome  risk  requirement  component  principle  constraint  pdr  adr  cr
Parts:
  criterion  edge-case  option  consequence
Relationships:
  mitigation  dependency  fulfillment  collaboration  effect
```

That recovers the grouping the nested form gave for free, without putting it in the syntax. Per-word
help is unaffected — `intent add mitigation --help` still states exactly what it connects.

**One rule the flat form deletes.** Element keys never enter the command namespace: they are always
positionals or flag values. So there is no reserved-word problem — a component keyed `edge-case` is
fine.

### Verb × thing matrix

| | add | rm | set | show | list | tree | trace | impact | promote/demote |
|---|---|---|---|---|---|---|---|---|---|
| product | — | — | ✓ | ✓ | — | ✓ | ✓ | ✓ | — |
| architecture | — | — | ✓ | ✓ | — | ✓ | ✓ | ✓ | — |
| job | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | — |
| outcome | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ |
| risk | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | — |
| requirement | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ |
| component | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ |
| principle | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | — |
| constraint | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | — |
| pdr / adr / cr | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | — |
| *parts* | ✓ | ✓ | — | — | ✓ | — | — | — | — |
| *relationships* | ✓ | ✓ | — | — | ✓ | — | — | — | — |

Roots cannot be added or removed. Parts and relationships have no identity of their own, so they
support only `add`, `rm`, and `list`.

### `architecture`

`architecture` is the top-level entity under the engineering domain — the engineering root element,
which renders as the architecture document.

- **Kind:** `KindEngineering` → `KindArchitecture` (`internal/tree/tree.go:27`).
- **Domain:** unchanged. `DomainEngineering` stays `engineering`, and so does `--domain` on `list`.
- **Storage:** unchanged. `Root.Engineering` and the `engineering:` YAML key stay as they are
  (`internal/model/model.go:12`).

This makes `Kind != Domain` for engineering, where they coincide for product — one translation point
that needs a test: `effect` may target a **root literal** (`internal/cli/link.go:36`), so
`--to architecture` must be written to storage as `engineering`. It is the only place a CLI word maps
to a different storage word.

## 4. Elements: flags

Derived from `internal/model/model.go`. Scalar fields become `set`/`add` flags; list fields become
parts (§5); relationships get creation-time seed flags (§6).

| element | scalar flags | part seeds (repeatable) | relationship seeds (repeatable) | parent |
|---|---|---|---|---|
| product | `--name --summary --detail` | — | — | — |
| architecture | `--name --summary --detail` | — | — | — |
| job | `--name --story --detail` | — | — | — |
| outcome | `--name --statement --detail --document` | — | — | `--job` |
| risk | `--name --statement --detail` | — | — | `--outcome` |
| requirement | `--name --statement --detail --document` | `--criterion` | `--mitigation` (req), `--dependency` | `--outcome` |
| component | `--name --responsibility --data-model --interfaces --behavior --detail --document` | `--edge-case` | `--fulfillment` (req), `--collaboration` | — |
| principle | `--name --statement` | — | — | — |
| constraint | `--name --statement` | — | — | — |
| pdr / adr | `--name --summary --context --decision` | `--option`, `--consequence` | `--effect` (req) | — |
| cr | `--name --change --rationale` | — | `--effect` (req) | — |

`(req)` marks a relationship required at creation, because the tree is validated before it is written
(`internal/cli/add.go:38-42`) and cannot hold a requirement without a `mitigates`, a component without
a `fulfills`, or a record without an `affects`. These seed flags are the one deliberate duplication
with the relationship words: **seed at creation, maintain with the relationship word.** Their names
are the relationship words, so there is one vocabulary rather than two.

`set` accepts the same scalar flags, plus positional-change flags. It never touches parts or
relationships.

```
set <element> <loc> --key <new-key>          # rename     (was mv --rename)
set <element> <loc> --<parent> <loc>         # re-parent  (was mv <addr> <parent>)
set <element> <loc> --before|--after <sib>   # reorder    (was mv --before/--after)
```

## 5. Parts

A part is a **repeatable component of one element**. It has no key; the locator flag names its owner.

| part | owner | verbs | source |
|---|---|---|---|
| `criterion` | requirement | `add` `rm` `list` | `acceptance_criteria` (authored) |
| `criterion` | component | `list` only | **derived** — union of `acceptance_criteria` across `fulfills` |
| `edge-case` | component | `add` `rm` `list` | `edge_cases` (authored) |
| `option` | pdr / adr | `add` `rm` `list` | `options` (authored) |
| `consequence` | pdr / adr | `add` `rm` `list` | `consequences` (authored) |

```
intent add  criterion "Token expires in 15m"    --requirement token-ttl
intent list criterion                           --requirement token-ttl
intent rm   criterion --index 2                 --requirement token-ttl
intent add  edge-case "Upstream 503 mid-retry"  --component api-gateway

intent list criterion --component api-gateway   # derived; no add/rm
```

`criterion` is the one word whose owner flag also selects its meaning: `--requirement` reaches an
authored list, `--component` reaches a computed one. `add criterion --component X` must be refused
with a message naming the derived view, not a generic flag error.

**Selectors for `rm`.** `--index N` is 1-based against the order `list` prints — deterministic, and
the number is already on screen. `--match <substring>` is a convenience that errors on zero or
multiple matches. Exactly one of the two is required.

### `success_criteria` is dropped from the schema

A component's reason to exist is its `fulfills` edges — *"every component should earn its place by
fulfilling a requirement"* (`internal/help/content/concept/component.md:22`). So if a component
satisfies every requirement it fulfills, and those requirements carry acceptance criteria, the
component's success is already fully specified upstream. That leaves `success_criteria` only two
possible contents:

1. **A restatement of the acceptance criteria of the requirements it fulfills** — duplication, which
   the trace graph exists to prevent.
2. **A condition nobody wrote as a requirement** — which breaks the invariant that every condition
   traces to a risk. A requirement smuggled into the engineering domain.

A third strike comes from the framework's own judgment layer: a component is *"a charter, not an
outcome"* (`concept/component.md:26`). Success criteria measure achievement, which is what `outcome`
is for, so the field is an altitude violation baked into the schema — `judgment:altitude` argues
against it.

The escape hatch already exists. A component-level condition that genuinely isn't a product
requirement ("the renderer must be deterministic") is an `engineering.constraints` entry, which is a
keyed map precisely so it carries a stable address and can be an edge target
(`internal/model/model.go:28-29`).

So the component view becomes **derived, not authored**. It cannot drift, and it matches the tool's
own premise — *"The YAML is canonical; generated docs are derived"* (`internal/cli/root.go:31`).

One honest caveat on the evidence: `success_criteria` appears zero times in this repo's own
`intent.yaml` — but so do `edge_cases`, `data_model`, and `interfaces`. Four of the six component
structured fields are unused, so the count says the components here are thin, not that the field is
wrong. The argument above is the case; the usage count is only consistent with it.

Touches `Component.SuccessCriteria` (`model.go:103`), its `appendList` in the index walk
(`tree.go:242`), `/$defs/component` in `internal/schema/intent.schema.json`, and
`ref/component.md:19` / `concept/component.md:17`. No migration tooling at v0 — hand-edit.

## 6. Relationships

`link`/`unlink` (`internal/cli/link.go`) dissolve into five words. Each word determines **both**
endpoint kinds, so the flags are derivable and each word's help states exactly what it connects.

| word | connects | flags | storage |
|---|---|---|---|
| `mitigation` | requirement → risk | `--requirement` `--risk` | `Requirement.Mitigates` |
| `dependency` | requirement → requirement | `--requirement` `--to` | `Requirement.DependsOn` |
| `fulfillment` | component → requirement | `--component` `--requirement` | `Component.Fulfills` |
| `collaboration` | component → component | `--component` `--to` `--note` | `Component.Collaborations` |
| `effect` | pdr\|adr\|cr → any element | `--pdr`\|`--adr`\|`--cr` `--to` | `*.Affects` |

Per R3: the source is always kind-named; the target is kind-named when its kind is fixed and differs
from the source, and `--to` otherwise. `--to` on `effect` is repeatable — a record usually touches
several things.

```
intent add  mitigation    --requirement token-ttl --risk cart-abandon
intent add  dependency    --requirement token-ttl --to session-store
intent add  fulfillment   --component api-gateway --requirement token-ttl
intent add  collaboration --component api-gateway --to auth-service --note "calls to validate tokens"
intent add  effect        --cr fix-drift --to api-gateway
intent rm   fulfillment   --component api-gateway --requirement token-ttl

intent list mitigation  --risk cart-abandon        # what mitigates this risk
intent list fulfillment --requirement token-ttl    # what fulfills this requirement
intent list dependency  --to session-store         # what depends on this
```

Because the source flag is kind-named, argument order is free and the read path works from either
end without knowing the direction convention. That is why kind-named flags beat a uniform
`--from`/`--to`.

### `relationships:` → `collaborations:`

The component-to-component field is renamed in storage as well as in the CLI. Two reasons:

1. **`relationships` is a category naming one of its members.** All five words above are
   relationships; only one of them was called that.
2. **`Responsibility` + `Collaborations` completes the CRC card** — Class, Responsibility,
   Collaboration. The element key and `name` are the Class slot, `responsibility` is already there,
   and this field is the third. No new field is needed; document the mapping in
   `concept/component.md`.

`coupling` was the earlier candidate and is wrong for this field: it implies structural dependence
and reads as a smell. The four real notes in this repo describe runtime collaboration —
*"never renders a tree that failed validate"*, *"every write re-validates through this component
before committing"*.

Touches `Component.Relationships` (`model.go:98`), `EdgeRelationships` → `EdgeCollaboration`
(`tree.go:56`), the `addEdge`/`removeEdge` arms (`internal/mutate/edges.go:26,60`),
`/$defs/component`, `ref/edges.md:15`, and 4 occurrences in `intent.yaml` (lines 214, 225, 236, 250).

**A narrowing to note.** `ref/edges.md:15` documents this edge as pointing at "another element",
but `addEdge` never type-checks the target and all four real uses point at components. The flat form
requires each word to fix both endpoint kinds, so this is specced as **component → component**. That
is a deliberate tightening of the documented contract, consistent with every existing use.

### Two things this tightens generally

`addEdge` (`internal/mutate/edges.go:14`) does not type-check edge targets today; only E001
(dangling) and E004 (cross-domain record) catch mistakes after the fact. Declaring endpoint kinds
per word means the resolver rejects a bad target at parse time — a real behavior change and a strict
improvement.

And `mitigates` points at **same-outcome** risks (`internal/model/model.go:77-78`), so
`add mitigation --requirement X --risk Y` can scope `Y` to `X`'s own outcome subtree, making bare
risk keys work without qualification. Derive the scope; don't ask for it.

### `effect` and `impact`

Both survive, and the adjacency is deliberate: `add effect` declares an edge, `impact` reads the
derived closure over all edges. That is the same authored-vs-derived pairing already accepted for
`criterion` (§5), and the two never appear in the same argument position.

## 7. Locator resolution

An additive change to `internal/tree/address.go`. `Resolve` becomes a thin wrapper over:

```go
type ResolveOpts struct {
    Kind  tree.Kind // "" = any
    Scope string    // "" = whole tree; else an ancestor address
}

func (ix *Index) ResolveIn(query string, opts ResolveOpts) (*Element, error)
```

1. If `query` is an exact full address in `byAddr` and satisfies `Kind` and `Scope`, return it.
   Exact addresses always win (today's `address.go:39`).
2. Otherwise build candidates from `ix.order`, filtered by `Kind` and by `Scope` prefix.
3. Match `segmentSuffix` (`address.go:69`) against that candidate set.
4. Zero → `NotFoundError`. One → the element. More → `AmbiguousError` with full addresses.

`ResolveIn(q, ResolveOpts{})` is byte-for-byte today's behavior, so the change is purely additive.
Errors must name the scope, or the new precision is invisible:

```
no requirement matches "auth"

"auth" is ambiguous among requirements; candidates:
  product.jobs.checkout.outcomes.fast.requirements.auth_token
  product.jobs.signup.outcomes.secure.requirements.auth_flow
```

Scope comes from the locator flag — `--outcome fast-checkout` scopes a `--risk` locator to that
outcome's subtree — or from a derived constraint, like same-outcome mitigation (§6).

## 8. Migration

| old | new |
|---|---|
| `tree` | `tree` |
| `show <addr>` | `show <element> <loc>` |
| `find [q] --type T --domain D` | `list [<thing>] [q] --domain D` (`--type` dropped; the thing is the filter) |
| `trace <addr> --edges` | `trace <element> <loc> --via` |
| `affects <addr>` | `impact <element> <loc>` — now transitive, §2 |
| `validate` | `validate` |
| `build --out` | `build --out` |
| `check --out` | `build --check --out` |
| `add <type> <parent> <key> …` | `add <element> <key> [--<parent> <loc>] …` |
| `set <addr> <field> <value>` | `set <element> <loc> --<field> <value> …` (multiple fields per call) |
| `set <addr> type document` | `promote <element> <loc>` |
| `set <addr> type inline [--cascade]` | `demote <element> <loc> [--cascade]` |
| `rm <addr>` | `rm <element> <loc>` |
| `link mitigates <req> <risk>` | `add mitigation --requirement <loc> --risk <loc>` |
| `link dependsOn <req> <req>` | `add dependency --requirement <loc> --to <loc>` |
| `link fulfills <comp> <req>` | `add fulfillment --component <loc> --requirement <loc>` |
| `link relationships <comp> <comp> --note` | `add collaboration --component <loc> --to <loc> --note` |
| `link affects <record> <addr>` | `add effect --pdr\|--adr\|--cr <loc> --to <loc>` |
| `unlink <kind> <from> <to>` | `rm <relationship-word> --<source> <loc> --<target> <loc>` |
| `promote <addr>` | `promote <element> <loc>` |
| `mv <addr> [parent] --rename --before --after` | `set <element> <loc> [--<parent> <loc>] [--key <new>] [--before\|--after <sib>]` |
| `record <pdr\|adr\|cr> <key> …` | `add pdr\|adr\|cr <key> …` |
| `install-skill --agent --dir` | `install skill --agent --dir` |
| `help [slug] --list --json` | unchanged |

`--file`/`-f` stays a persistent root flag.

One capability is lost and must be preserved deliberately: `unlink` can remove a **dangling** edge
whose target no longer resolves, via `resolveEdgeTarget` (`internal/cli/link.go:73-85`). `rm
<relationship-word>` needs the same fallback — match the query against the addresses the source
actually declares — or hand-edited breakage becomes unfixable through the CLI.

### Cost

The cobra rewiring is the small part.

- **`internal/help/content/` — 41 topic files.** `ref/commands.md`, `ref/addressing.md`,
  `ref/edges.md`, `ref/syntax.md`, all five `guide/*`, all five `errors/*`, and every `concept/*` and
  `ref/*` page that shows a command. This is the bulk of the work.
- **The generated agent skill.** `install-skill` renders the skill from that embedded help
  (`internal/cli/root.go:57`), so every agent consuming `intent` gets the new grammar only after the
  help rewrite lands. The help *is* the skill.
- **Dogfooded docs.** This repo builds `docs/product/` and `docs/engineering/` from its own
  `intent.yaml`; both need a rebuild, and `intent.yaml`'s own prose describes the old commands.
- **Schema.** Two storage changes: drop `success_criteria` (§5), rename `relationships:` →
  `collaborations:` (§6). No migration tooling at v0 — hand-edit this repo's `intent.yaml` (4 lines
  for the rename, 0 for the drop).
- **Tests.** `internal/cli/*_test.go` is written against the flat surface throughout.

## 9. Settled

| # | question | decision |
|---|---|---|
| 1 | subnoun shape | **Flat everywhere.** One slot after the verb; `rel` dissolves; taxonomy moves to grouped help (§3). |
| 2 | relationship endpoints | **Kind-named source, kind-named target where fixed and different, `--to` otherwise.** `--on`/`--with` dropped (§6). |
| 3 | the component-to-component field | **`collaborations:` / `collaboration`**, renamed in storage. Completes the CRC card; no `class` field — key + `name` is that slot (§6). |
| 4 | the `affects` verb | **`impact`, redefined as the transitive incoming closure.** Not a rename; a capability upgrade (§2). |

### Small calls made along the way, flagged rather than asked

- `trace --edges` → `--via`, since "edge" left the user-facing vocabulary (§2).
- `collaboration` is specced **component → component**, narrowing `ref/edges.md:15`'s "another
  element" to match all four existing uses (§6).
- `impact` gets a dedicated `Index.Impact(addr)` rather than overloading `Trace`'s depth clamp (§2).
