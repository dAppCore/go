---
title: CoreGO Development
description: How to build, test, and hold CoreGO to its own gold-standard quality bar.
---

# CoreGO Development

## Build & verify

```bash
core go test          # run the suite
core go qa            # fmt + vet + lint + test
core go qa full       # + race, vuln, security
```

CoreGO is the root module `dappco.re/go` with **no `external/` submodules** — it
builds against the Go standard library alone. The Go 1.26 workspace lives at
`~/Code/go.work`; after adding modules run `go work sync`.

## The test standard

Every exported symbol carries one focused test **per variant**:

```
Test<File>_<Symbol>_{Good,Bad,Ugly}        // function
Test<File>_<Type>_<Method>_{Good,Bad,Ugly} // method (compound symbol)
```

Rules, enforced by the audit (below) and reviewed by hand:

- **One test per symbol per variant** — no duplicates; the name mirrors the real
  code symbol, never a shorthand or scenario.
- **No gamed tests** (`AssertTrue(t, true)` shells) and **no shallow tests** (a
  lone happy path). Good/Bad/Ugly pin the contract across distinct cases —
  value, edge, and error/contention paths. Prefer no test over a gamed one.
- **`X_test.go` is the test file** per `X.go` (white-box `package core` where a
  test needs unexported access; black-box `package core_test` otherwise). The
  file set is `X_test.go` + `X_example_test.go` + `X_bench_test.go`.
- Test-only helpers that must reach unexported code live in `package core` test
  files (e.g. `assertStub`, `ExecCmdForTest`); black-box tests see their
  exported helpers through the test-augmented package.

Helpers: assertions are `core.Assert*`/`Require*` (they take `testing.TB`, so a
`*stubT` embedding `*T` can capture failures); subprocess tests re-exec via
`ExecCmdForTest(Args()[0], "-test.run=^TestX$")` gated by an env var.

## The audit

`cga.sh` (CoreGoAudit, in the `lethean-claude` skill) is how names and coverage
stay honest at ecosystem scale — CoreGO must pass the audit it holds every
consumer to. Key buckets and their scripts:

| Bucket | Means |
|--------|-------|
| `unreferenced-tests` | a test body never names its target symbol (dispatcher gaming) |
| `ax7-triplet-gaps` | a public symbol missing a Good/Bad/Ugly |
| `identical-triplets` | Good/Bad/Ugly with byte-identical bodies (three copies, not three cases) |
| `test-stubs` | `Test*` with ≤2 body lines (a *signal* — a concise real assertion is fine; a no-op is not) |
| `example-gaps` | a public symbol with no `Example<Symbol>` |
| `unbenched` | a symbol no benchmark exercises (the bench+coverage method) |

`test-stubs` is a heuristic, not a mandate: `AssertEqual(t, "42", Itoa(42))` is a
complete one-line test and must **not** be padded to satisfy a line count —
that is itself gamed work.

## Contributing

- UK English; EUPL-1.2.
- Errors via `core.E(...)`; no banned stdlib imports (see
  [architecture.md](architecture.md)).
- API-shape changes ripple to all `dappco.re/*` consumers — coordinate them as a
  release, don't make them in isolation.
- Commit trailer: `Co-Authored-By: Virgil <virgil@lethean.io>`.
