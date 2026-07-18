<!-- SPDX-License-Identifier: EUPL-1.2 -->
# core/go v0.12.0 — the named-resource law + the sealed bundle

**Theme:** every subsystem speaks one grammar — `c.Subsystem(name).Verb(opts) → Result → .Typed()`.
Named resource in, uniform handle out, typed value out the far side. A `core.New(...)` bundle,
once locked, is an immutable capability toolkit exportable as a package var.

Tags: **`decide:`** implementer must choose (don't guess) · **`verify:`** confirm before acting ·
**`breaking:`** consumer-visible shape change (this window only) · **`additive:`** lands any time.

Evidence base: fleet campaign 2026-07-18 (fences A–D, commits `4c702a2`/`879a496`/`a520b60`/`b041db5`),
cga baseline 496 findings pre-fleet, adjudication verdict 446 stub flags → 6 genuine (98.7% honest).

---

## Wave 1 — additive (landing first, nothing breaks)

| # | Item | Shape | Status |
|---|---|---|---|
| W1-1 | **Result typed getters** — mirror the Options accessor dialect on the universal output: `Int() Bool() Float64() Duration() Bytes() Err()`, strict zero-value-on-miss contracts identical to options.go (Float64 promotes int/int64/float32; Duration accepts Duration or ParseDuration string). `String()` is the ONE deliberate divergence: it is Stringer-first (value rendered when OK, Error() text when not) so `%v` of any Result reads well in logs — documented loudly at the method. | additive | this session |
| W1-2 | **`Return[T]`** — the typed twin, NOT a second universal: `{Value T; Err error}`, flat, Err is the discriminant. Constructors `ReturnOK/ReturnFail/ReturnFrom(v,err)/ReturnOf[T](Result)/ReturnTry[T](fn)` (panic→Err). Methods `OK() Or(T) T Must() T Code() Error() Result() Log(op,msg)` — Log delegates to the package log path when Err != nil and returns self for chaining, so consumers never hand-wire logging; ambient auto-log rejected (expected-failure branches like fs.notfound would spam). | additive | this session |
| W1-3 | **The three-line law** (RFC + AGENTS.md): (1) anything crossing a type-erasing boundary — IPC, registries, Actions, Tasks, RegistryOf — returns `Result`, no exceptions; (2) leaf functions with one concrete struct-shaped return may declare `Return[T]`; (3) scalars use Result + getters, never lifted into Return. One symbol never offers both. | docs | this session |
| W1-4 | **`c.API(name)`** — the named-resource accessor for remotes: `API(name ...string)` zero-arg returns the subsystem (unchanged); named returns a bound view (`On(name)`) whose `Invoke(action, opts)` and zero-arg `Stream()` use the binding. The resolution chain already existed end-to-end (Stream→Drive→scheme→protocol); this is accessor sugar only. `Stream(name string)` → `Stream(name ...string)` (verify: pragmatically compatible; method-value assignments would break — none found in-repo). | additive | this session |
| W1-5 | **`MustNew`** — `New()` on a failed option logs and returns a half-built Core (contract.go); fatal for the package-var bundle pattern. `MustNew` panics on the first failed option: fail at import, correct for `var Bundle = core.MustNew(...)`. New keeps its forgiving shape for main(). | additive | this session |

## Wave 2 — the breaking window (the release's reason to be v0.12.0)

**Status 2026-07-18 (same session):** W2-1 LANDED (cli.noop/cli.unknown sentinels; a
typo'd command now exits non-zero; RunResult converts only the sentinel). W2-2 LANDED
for Service (coded miss `core.Service: service not found`); Cli()/full one-miss-law
table remains open. W2-3 LANDED (sealConclave at ServiceStartup tail — Locked:
actions/tasks/commands/protocols/drive/data; Sealed: features; open: locks/config).
W2-4 LANDED (merge). W2-5 was ALREADY DONE — Setenv/Unsetenv have returned Result
since v0.9; the audit's 2 hits were its own doc comments (instrument fixed). W2-6
instrument-resolved: the return-position regex fix + wrapper-mode exemptions leave 0
genuine unconverted funcs. W2-7 LANDED (all 6 discards now handled explicitly; the
embed close-after-failed-write discard was CORRECT and is now a commented branch).
W2-8 decided: keep `Call(endpoint, action)` as the primitive and `Invoke(action)` as
the bound verb — no collapse.

| # | Item | Fix |
|---|---|---|
| W2-1 | **RunResult nil-failure conflation** (core.go ~205): `!r.OK && r.Value == nil → OK:true` silently converts any valueless failure into success | CLI's "no commands, banner shown" no-op returns a sentinel code (`cli.noop`); RunResult checks `r.Code() == "cli.noop"`, never infers from nil. breaking: handlers relying on the swallow (none should) surface. |
| W2-2 | **The one-miss law** — today: Action→zombie, Service→empty `Result{}`, RegistryOf→coded error, Lock→get-or-create, Cli→typed nil | decide + document per category: resources (Lock, API(name), Drive) get-or-create; capabilities (Action, Task, Command) zombie-with-Exists; queries (Service, RegistryOf) coded-error Result. Service's bare `Result{}` gains a code; Cli() returns coded-error via ServiceFor miss (or keeps nil + nil-receiver guards on every *Cli method — decide). |
| W2-3 | **Conclave-wide sealing** — WithServiceLock freezes ONLY services (lock.go:76-89); the "immutable bundle" is porous | at LockApply: **Locked** (frozen) services, actions, tasks, commands, protocols, drive, data; **Sealed** (no new keys, values mutable) features + config — flags stay toggleable, none inventable. Registry modes already exist; wire them. breaking: post-lock RegisterProtocol/Action registration fails where it silently succeeded. |
| W2-4 | **WithOptions clobbers** (contract.go:162) — replaces the whole options struct, discarding earlier WithOption keys | merge semantics (Set per item). breaking only for anyone depending on the clobber. |
| W2-5 | **`Setenv`/`Unsetenv` → Result** (cga breaking-api-sites: 2) | the last two audit-flagged raw-error sites in the wrapper surface. |
| W2-6 | **err-shape residue** (cga: 10 `func … error` in production) | convert to Result. exception (H1, unchanged): error *constructors* (E/Wrap/WrapCode/NewCode/NewError/Errorf) keep returning error — they build the value a Result wraps. |
| W2-7 | **result-discards** (cga: 6 `_ =` in production) | handle or annotate each; zero silent drops. |
| W2-8 | **Call signature collapse** — `API.Call(endpoint, action, opts)` + bound `Invoke(action, opts)` coexist after W1-4 | decide: fold to `Call(action, opts)` on bound views only, endpoint form via `c.API(name)`; or keep both. One verb should win the docs. |

## Wave 3 — wiring the dormant features (additive, each needs one decision)

| # | Item | Decision needed |
|---|---|---|
| W3-1 | **RemoteAction into dispatch** — the colon law. `c.RemoteAction` fully resolves `"host:action"` (api.go:159) but has zero production callers | decide: wire into `Action.Run` when name contains ":" (transparent, matches the godoc vision) vs keep explicit. Transparent makes `c.API(name)` + Action registry one namespace — the vision — but adds a hot-path branch (bench it; Action.Run is per-dispatch). |
| W3-2 | **PerformAsync taskID injection** — handlers can't call `c.Progress(taskID,…)` because they never learn their ID (action.go:284-299) | inject as reserved `_task` key. verify: Options copy semantics — PerformAsync receives Options by value but the backing items slice is shared; Set on it may leak to the caller. Clone-then-Set. Also decide: dispatch on `c.context` (shutdown-aware) instead of `Background()`. |
| W3-3 | **RecordUsage at the choke point** — Entitled(action, quantity) + UsageRecorder both exist; Action.Run checks Entitled but never records on success | close the metering loop in Action.Run: entitled → run → OK → RecordUsage(name). decide: quantity source (Schema? opts key? default 1). |
| W3-4 | **Action.Schema validation** — declared (action.go:43), never enforced | opt-in: validate opts keys against Schema when Schema is non-empty; Result code `action.schema` on miss. Bench the non-empty path. |
| W3-5 | **OnReload runner** — field exists, no runner, no trigger (lifecycle.md documents the gap honestly) | decide: `ServiceReload(ctx)` runner + `Reloadable` interface discovery in RegisterService; trigger stays consumer-called (no config-watch until Config grows one). Or delete the field — half a feature is worse than none. |
| W3-6 | **lsp.go anchoring** — 25K, biggest file in the repo, no Core accessor, absent from the subsystem table | decide: `c.LSP()` accessor + CLAUDE.md table row, or extract to `dappco.re/go/lsp` consumer package. Its registries are already dogfooded (A7). |
| W3-7 | **`--- Global Instance ---`** dangling section header at core.go EOF with nothing under it | delete the header or build the thing it promised. verify git history for intent. |

## Instrument fixes (skill repo — `lethean-claude/scripts`, NOT this repo)

The fleet adjudicated 446 stub flags: 6 genuine, 3 tool bugs, 2 false-positive classes. The residual
cga count is analyzer debt, not repo debt — fix the instrument before the next campaign:

1. `test-stubs.py` brace counter breaks on `{`/`}` inside string/rune literals (TestJson_JSONDelim, body counted as 0).
2. `test-stubs.py` `func Test` regex matches inside raw-string fixtures (lsp_bench_test.go's embedded test-content).
3. `example-gaps.py` demands `Example<Symbol>` spellings Go's vet grammar forbids for underscore-bearing identifiers (`SHA3_256`, `O_NOFOLLOW`) — special-case to the package-level `Example_x` form.
4. action-name-format greps doc comments (api.go's two "findings" are godoc text).
5. service-name-empty flags the deliberate Bad-path test asserting the empty-name error.
6. decide: stub threshold (≤2 body lines) fires on honest exact-value asserts at 98.7% — gate on assertion-presence/tautology instead of line count, or accept the dimension as advisory on wrapper-layer repos.

## Release mechanics

- Wave 1 tags a v0.11.x (additive). Waves 2–3 are the v0.12.0 window proper.
- Consumer cascade per the corego-release skill (dev→main merge order, pinned tags, api↔inference cycle).
- Docs regenerate WITH the release: knowledge-pack refresh from `go doc -all` + example tests
  (the pack's invented-API drift is what masked these gaps — see memory `corego-knowledge-pack-drift`).

---
*Wave 1 lands 2026-07-18 (this session). Baseline: 496 cga findings pre-fleet; post-fleet number in
the module-gate receipt. The 6 genuine test fakes are fixed in the fence commits; everything still
flagged after instrument fixes is real wave-2/3 work.*
