---
title: CoreGO Architecture
description: How the CoreGO framework is composed — the Core container, the universal primitives, and the AX design rules.
---

# CoreGO Architecture

CoreGO (`dappco.re/go`) is the **L0 root** of the Core ecosystem: every other
`dappco.re/*` module depends on it, and it depends on nothing outside the Go
standard library. It has no `external/` submodules — it *is* the substrate.

## The Core container

Everything hangs off a single `*Core`, built with functional options:

```go
c := core.New(
    core.WithOption("name", "agent-workbench"),
    core.WithService(cache.Register),
    core.WithServiceLock(),
)
```

`New` applies each option in order, then seals the service registry
(`LockApply`). A `*Core` owns the subsystems — `App`, `Config`, `Data`, `Drive`,
`Fs`, `I18n`, `API`, the service/command/action/task registries, the error and
log panics, and the entitlement checker — and wires itself into each.

## Universal primitives

| Primitive | Role |
|-----------|------|
| `Result` | The single return shape. `{Value any, OK bool}` — value on success, error in `Value` on failure. Replaces `(T, error)` tuples. |
| `Registry[T]` | The universal brick: a thread-safe, ordered, lockable/sealable named collection. Services, commands, actions, locks, protocols and feature flags are all `Registry[T]`. |
| `Action` | A named callable with panic recovery + an entitlement gate. `Task` composes Actions. |
| `Service` | A lifecycle component (`OnStart`/`OnStop` return `Result`), discovered or named, registered on the Core. |
| `Command` | A path-based CLI tree (`deploy/to/homelab`), parents auto-created. |
| `Feature` | A handle over the feature-flag registry (`c.Feature("dark").Enable()`). |

`Result` is the spine: I/O, actions, commands, services and the assertion
helpers all speak it, so call sites compose uniformly (`if !r.OK { return r }`).

## Concurrency & lifecycle

`Core.Go` runs a function on a tracked `WaitGroup`; long-running goroutines poll
`Core.IsShutdown` to exit cooperatively. `ServiceShutdown` sets the shutdown
flag, cancels the root context, drains the WaitGroup, then stops services in
reverse dependency order. `WithServiceLock` seals registration after
construction so the running set is immutable.

## The AX design rules (why the code looks the way it does)

CoreGO is written to be navigable by agents as much as humans (RFC-CORE-008):

- **Universal types over bespoke returns** — `Result`/`Registry[T]` everywhere.
- **Predictable names + path-is-documentation** — `X.go` ↔ `X_test.go` ↔
  `X_example_test.go` ↔ `X_bench_test.go`; test names mirror the code symbol
  (`TestFile_Type_Method_{Good,Bad,Ugly}`).
- **Banned imports in consumer code** — no direct `os`, `os/exec`, `fmt`, `log`,
  `errors`, `strings`, `path/filepath`, `encoding/json`. CoreGO re-exports these
  through its own typed, `Result`-shaped wrappers (`core.Open`, `core.E`,
  `core.Contains`, `core.JSONMarshal`…) so the whole ecosystem has one surface.
- **Errors via `core.E(op, msg, cause)`** — structured, wrapping, `Is`-traversable.

## Module layering

CoreGO is L0. Higher modules (`dappco.re/go-process`, `…/go-inference`, the
agent/CLI layers, the LEM engines) consume it and register services/actions
into a `*Core`. Because CoreGO carries no third-party dependencies, a change
here ripples to all consumers — API-shape changes are release-coordinated, not
made in isolation.

See [RFC.md](RFC.md) for the authoritative contract and
[development.md](development.md) for the build/test/audit workflow.
