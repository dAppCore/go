<!-- SPDX-License-Identifier: EUPL-1.2 -->
# core/go — issues

## Progress (suite green at 3547 tests throughout)
**All hand-found items done (21):** J1 (cli opt-in tests) · I1 (Core cli comment) ·
F1 (Registry.Get/Has gate disabled + `GetIncludingDisabled` + dispatch test) ·
H2 (`fs.TempDir`→Result, wraps error in op context) ·
H1 (`WalkDir`/`HTTPListenAndServe`/`IsDir`/`IsFile`/`Exists`→Result, ~40 call sites) ·
E1 (RandPick misuse-panic: documented + AssertPanics-tested) ·
E2 (Once.Reset not-during-use contract: documented + tested) ·
A1 (Lock → Core-held `Registry[*Lock]`; added `Registry.GetOrSet`; Lock now pure DTO) ·
A2 (feature flags → `Registry[bool]`; actions keep `a.enabled` — hot-path + F1, documented) ·
A3 (Command child map deleted; flat registry IS the tree) ·
A4 (embed `assetGroups`/`AssetGroup.assets` → `Registry`, mutexes gone) ·
A5 (i18n: kept slice [no key] + private mu [encapsulated guard], documented) ·
A6 (SysInfo.values: read-only-after-init → left) ·
A7 (LSP `documents` → `Registry[[]byte]`; Suffixes left a value set) ·
A8 (ipcRegMu/queryRegMu/log.mu/crashMu: investigated → all left, with rationale) ·
B1+B2+B3 (`RegistryOf`→`Result`, 10 registries, snapshot-documented) ·
C1 (`LockEnable`/`LockApply` drop vestigial name → service-global) ·
G1 (`c.Config(group)` prefix view + `c.Feature(name)` handle) ·
D1 (example files for format/iter/lsp/unsafe/user).

**Remaining (mechanical, out of the inline scope):** K-series (cga static sweep) · L-series (benchmarking).


**Goal:** core/go should use its own primitives — `Registry[T]` for keyed collections, `c.Lock(name)`
for locks, `Result` for returns — not hand-rolled maps/mutexes/raw returns in-package.

Tags in the Fix column: **`decide:`** a choice the implementer must make (don't guess) · **`verify:`**
confirm the claim before acting · **`depends on:`** ordering constraint · **`exception:`** a carve-out.

## Registry collections — hand-rolled `collection + private-mutex` → `Registry[T]`
`Registry[T]` already bundles a guarded map, so adopting it removes the private mutex too.

| # | Problem | Fix |
|---|---|---|
| A1 | Lock's `locks SyncMap` is embedded in the `Lock` DTO itself (lock.go:14-39) | Core holds `locks *Registry[*Lock]`; `c.Lock(name)` becomes Get-or-Set on it; `Lock` drops the `locks` field → pure DTO (`Name`, `Mutex`). |
| A2 | Config feature flags (`Features map[string]bool` + `mu`, config.go:58/78) and Action enablement each re-implement `Enable`/`Disable`/`Enabled` (config.go:159-197, action.go:85-108) — logic `Registry` already has | back feature flags with a `Registry[bool]` and actions with the actions registry's enable/disable; keep the public `Enable/Disable/Enabled` methods as one-line shims over it. |
| A3 | Per-`Command` child `commands map` (command.go:57) duplicates the flat path-keyed registry. verify: its only reads are the placeholder-preservation block (command.go:133-135); dispatch + help both use the flat `c.commands` (cli.go:83/138). | derive children via `c.commands.List("<path>/*")`; delete the child map **and** the preservation block. |
| A4 | embed package-global `assetGroups map[string]*AssetGroup` + `assetGroupsMu` (embed.go:50-51), and per-group `assets map[string]string` (embed.go:63) | `assetGroups` → `Registry[*AssetGroup]`; `assets` → `Registry[string]`; both mutexes go. |
| A5 | i18n uses a private `mu` (i18n.go:53) + a hand-managed `locales []*Embed` (i18n.go:54) | mutex → `c.Lock("i18n")`. decide: `locales` is an ordered append-list with **no natural key** — either keep it a slice (just fix the lock) **or** introduce a key (mount path) to make it `Registry[*Embed]`. Do not force a registry without a key. |
| A6 | SysInfo `values map[string]string` (info.go:87) | verify: is `values` mutated after init, or read-only system info? If mutated/concurrent → `Registry[string]`; if read-only → leave it (not a registry case). |
| A7 | LSP `documents map[string][]byte` + `docsMu` (lsp.go:187-188), `Suffixes map[string]bool` (lsp.go:109) | verify each is keyed + mutable. `documents` → `Registry[[]byte]` (mutex goes); `Suffixes` → `Registry[bool]`, or leave as a value set if read-only. |
| A8 | Residual private mutexes: `ipc.ipcRegMu`/`queryRegMu` (ipc.go:21,24), `log.mu` (log.go:74), `crashMu` (error.go:476) | per mutex: if it guards a collection → fold into a `Registry`; else → `c.Lock(name)`. Investigate each — not a blanket change. |

## `RegistryOf()` (core.go:268-296)
| # | Problem | Fix |
|---|---|---|
| B1 | switch exposes only services/commands/actions; lock/data/drive/api/ipc-tasks/config/embed/i18n missing | add a case per registry. depends on A: a feature only appears here once it *is* a `Registry` (A1/A4/A5…). |
| B2 | unknown name → `NewRegistry[any]()`, indistinguishable from a genuinely empty registry (core.go:278) | return `Result` (OK=false on unknown name). breaking: changes the `*Registry[any]` return type — update callers. |
| B3 | `registryProxy` returns a copied snapshot, not a live view (core.go:283-…) | decide: return a live read-only view, **or** rename to `RegistrySnapshot()` and document the freeze. |

## `Lock` API
| # | Problem | Fix |
|---|---|---|
| C1 | `LockEnable`/`LockApply(name ...string)` accept a name then drop it — only a global `services.lockEnabled` (lock.go:79-91) | decide: implement per-name service locks (backed by the lock registry), **or** drop the `name ...string` param and document it as service-global. |

## `Registry.Disable` doesn't gate dispatch
| # | Problem | Fix |
|---|---|---|
| F1 | `Get`/`Has` ignore `disabled` (registry.go:86-105) but `List`/`Each` honour it (130/149); dispatch resolves via `Get` (cli.go:83) → a `Disable`d command/action still runs (the documented `Action("dangerous.purge").Disable()` doesn't stop it) | gate `Get`/`Has` (disabled → not-found); add `GetIncludingDisabled` for inspect/re-enable; audit existing `Get` callers for any relying on seeing disabled entries. test: `Disable("x")` → `Get` misses **and** `x` won't dispatch. |

## `Result` everywhere
| # | Problem | Fix |
|---|---|---|
| H1 | raw `error` returns (`WalkDir` fs.go:330, `HTTPListenAndServe` api.go:366) + naked `bool` predicates (`IsDir`/`IsFile`/`Exists` fs.go:380/400/420) | convert to `Result`. exception: the error *constructors* (`E`/`Wrap`/`NewCode`/`NewError`/`Errorf`, error.go) keep returning `error` — they build the value a `Result` wraps. |
| H2 | `fs.TempDir` returns `""` on `MkdirTemp` failure (fs.go:285-296) — swallows the error | return `Result`. |

## Config/Feature namespacing
| # | Problem | Fix |
|---|---|---|
| G1 | `c.Config()` is one flat bag — no group scoping, no first-class keyed feature accessor | depends on A2. add `c.Config("group")` → a group-scoped Config view, and `c.Feature(name)` → a keyed feature handle (name required). decide: the exact surface + which features opt into grouping ("not all, but quite a few"). |

## Tests / docs
| # | Problem | Fix |
|---|---|---|
| D1 | no `_example_test.go`: format.go, iter.go, lsp.go, unsafe.go, user.go | add runnable example tests on core/go's own `assert.go`. |
| I1 | Core struct comment (core.go:24) says cli is reached via `ServiceFor[*Cli]`, but `c.Cli()` is the real accessor (core.go:92) | point the comment at `c.Cli()`. |

## Consistency
| # | Problem | Fix |
|---|---|---|
| E1 | `RandPick` panics on empty slice (random.go:88-93) | decide: document the misuse-panic contract, or add a `Result`-returning variant. |
| E2 | `Once.Reset()` swaps the inner `sync.Once` (sync.go:164) — racy if another goroutine is mid-use | document the not-during-use contract, or guard the swap. |

## Test suite (regression from the WithCli change, commit 4b0072a)
| # | Problem | Fix |
|---|---|---|
| J1 | cli tests assume the old auto-registered cli — bare `New()` + `c.Cli()` across cli_test.go (~20 tests), core_example_test.go, command_test.go, core_test.go, example_test.go; suite panics → coverage unmeasurable | update each to `New(WithCli())`; then measure coverage. |

## Audit-surfaced tasks — `cga.sh` (the lethean-claude compliance audit)
A1–J1 are the hand-found framework-dogfood gaps. The **mechanical** house-standard sweep is the
lethean-claude CoreGoAudit — it reads the *whole* tree and emits the AX-7 / `Result`-idiom / test-shape
findings as counts. This is the same audit every consumer repo is held to; the dogfood goal is core/go
passing it.

```bash
bash ~/.claude/skills/lethean-claude/scripts/cga.sh /Users/snider/Code/core/go
# verdict + per-dimension counts; re-run after each batch to watch it fall.
# narrow while iterating:  cga.sh /Users/snider/Code/core/go ./pkg/foo/
```

Current verdict: **NON-COMPLIANT — 1488 findings.** Each dimension is a task bucket (work biggest-first):

| # | dimension | count | task |
|---|---|---|---|
| K1 | test-stubs | 901 | replace ≤2-line `Test*` shells with real assertions (the bulk of the debt) |
| K2 | example-gaps | 214 | add a runnable `Example<Symbol>` per exported symbol (extends D1) |
| K3 | unreferenced-tests | 140 | each test must name + exercise the symbol it claims — fake coverage otherwise |
| K4 | ax7-triplet-gaps | 109 | complete the `Test<File>_<Symbol>_{Good,Bad,Ugly}` triplet per symbol |
| K5 | action-name-format | 39 | rename actions to `domain.verb` |
| K6 | err-shape-funcs | 22 | `func … error` → `Result` — same work as H1, dedupe don't double-count |
| K7 | identical-triplets | 9 | three *cases*, not three copies of one assertion |
| K8 | missing-example-files | 9 | add `{file}_example_test.go` (subset of D1) |
| K9 | result-discards | 6 | stop dropping `Result` returns on the floor |
| K10 | tuple-result-shape | 6 | collapse `(T, error)` / `(T, bool)` → `Result` |
| K11 | missing-test-files | 3 | add `{file}_test.go` |
| K12 | docs-gaps | 2 | docs structure |
| K13 | service-name-empty | 1 | give the service a name |
| K14 | ax7-helpers | 2 | verify: audit-helper false positive vs a real gap |
| — | banned-imports (23) + breaking-api (2) | 25 | **core/go-definitional** — core/go *is* the wrapper layer that defines `Setenv`/`Unsetenv`/etc. verify case-by-case; the audit header still says it false-positives here (the skill-wording fix is its own task — don't "fix" core/go to satisfy a false positive). |

depends on J1: the audit's COVERAGE / DEAD-CODE / UNBENCHED sections currently report "build failed" because
the broken cli tests stop `go test ./...` compiling the suite. Fix J1 first → those three sections produce
numbers → they feed the benchmarking phase below.

## Benchmarking — the final phase (CRITICAL)
**CoreGO is the substrate under highly-optimised codebases (go-mlx and the LEM engines build on it); every
alloc and every byte-per-op there is paid on a hot, per-token path.** Benchmarks are not a nicety — they are
how the alloc/B-per-op floor gets *locked in* so a later edit can't silently regress it. This phase is last
because it can only measure once the suite builds (J1) and the surface is real (K-series), but for a framework
it is the one that matters most.

| # | task | how |
|---|---|---|
| L1 | every `{file}.go` ships `{file}_bench_test.go` | the third leg of the triplet (`_test` + `_example_test` + `_bench_test`); cga's UNBENCHED section is the gap list once it builds |
| L2 | bench the hot primitives with `-benchmem`; record allocs/op **and B/op** | `go test -run='^$' -bench=. -benchmem -benchtime=20x ./...` — read B/op as hard as allocs/op (use the `lethean-perf` skill's method) |
| L3 | drive each Registry / Lock / Result path to its alloc floor | profile (`-memprofile` → `pprof -alloc_objects -list`), prove the escape (`-gcflags=-m`); a non-empty map is a 2-alloc floor, an output buffer is 1 — evidence the floor, don't assert it |
| L4 | clear DEAD CODE (140 fake-coverage tests of 2562) | overlaps K3 — a test that never names its symbol is both fake coverage *and* dead weight in the bench/coverage build |

gate: L-series runs after J1 (suite builds) + the K-series (real tests exist to bench against). Until J1 lands,
cga's UNBENCHED reports "bench+cover build failed" — that's the J1 regression, not a benchmarking gap.

---
*Two task streams. **A–J** = hand-found framework-dogfood gaps (registries / locks / `Result`); ~22 of 65 files
swept by hand, the rest not yet read. **K–L** = the mechanical `cga.sh` sweep over the whole tree (1488
idiom/test findings + the benchmarking floor). Re-run `cga.sh` after each batch — the goal is a clean verdict:
core/go passing the audit it holds every consumer repo to.*
