---
title: Lifecycle
description: Run and RunResult entry points, startup, shutdown, context ownership, and background task draining.
---

# Lifecycle

CoreGO manages lifecycle through `core.Service` callbacks. `RegisterService` adapts the `Startable` and `Stoppable` interfaces onto those callbacks when an instance implements them.

## Entry Points: `Run` and `RunResult`

```go
c := core.New(core.WithService(myService.Register))
c.Run()
```

```go
r := c.RunResult()
if !r.OK {
	core.Exit(1)
}
```

### What `RunResult` Does

1. defers `ServiceShutdown` — it always runs, even on startup failure or panic
2. runs `ServiceStartup` on the Core context
3. runs the CLI when one is registered (`core.WithCli()` or a consumer CLI service)
4. returns the CLI's result

### Failure Behavior

- a failed startup returns that result; a non-error failure value is wrapped as code `core.run.startup`
- a CLI with nothing to run (no commands registered, banner shown) counts as success
- `Run()` is sugar over `RunResult()`: it logs the failure and calls `c.Exit(1)`, which runs the graceful shutdown chain

## Service Hooks

```go
c.Service("cache", core.Service{
	OnStart: func() core.Result {
		return core.Result{OK: true}
	},
	OnStop: func() core.Result {
		return core.Result{OK: true}
	},
})
```

Only services with `OnStart` appear in `Startables()`. Only services with `OnStop` appear in `Stoppables()`.

## `ServiceStartup`

```go
r := c.ServiceStartup(context.Background(), nil)
```

### What It Does

1. clears the shutdown flag
2. stores a new cancellable context on `c.Context()`
3. runs each `OnStart`
4. broadcasts `ActionServiceStartup{}`

### Failure Behavior

- if the input context is already cancelled, startup returns that error
- if any `OnStart` returns `OK:false`, startup stops immediately and returns that result

## `ServiceShutdown`

```go
r := c.ServiceShutdown(context.Background())
```

### What It Does

1. sets the shutdown flag
2. cancels `c.Context()`
3. broadcasts `ActionServiceShutdown{}`
4. waits for background tasks created by `PerformAsync`
5. runs each `OnStop`

### Failure Behavior

- if draining background tasks hits the shutdown context deadline, shutdown returns that context error
- when service stop hooks fail, CoreGO returns the first error it sees

## Ordering

`Startables()` and `Stoppables()` iterate the service registry in insertion order, so hooks run in registration order.

Shutdown also runs in registration order — it is not reversed. If teardown must happen in reverse dependency order, orchestrate it explicitly inside your `OnStop` callbacks.

## `c.Context()`

`ServiceStartup` creates the context returned by `c.Context()`.

Use it for background work that should stop when the application shuts down:

```go
c.Service("watcher", core.Service{
	OnStart: func() core.Result {
		go func(ctx context.Context) {
			<-ctx.Done()
		}(c.Context())
		return core.Result{OK: true}
	},
})
```

## Built-In Lifecycle Actions

You can listen for lifecycle state changes through the action bus.

```go
c.RegisterAction(func(_ *core.Core, msg core.Message) core.Result {
	switch msg.(type) {
	case core.ActionServiceStartup:
		core.Info("core startup completed")
	case core.ActionServiceShutdown:
		core.Info("core shutdown started")
	}
	return core.Result{OK: true}
})
```

## Background Task Draining

`ServiceShutdown` waits for the internal task waitgroup to finish before calling stop hooks.

This is what makes `PerformAsync` safe for long-running work that should complete before teardown.

## `OnReload`

`ServiceReload` runs every service's `OnReload` in registration order — trigger it from a signal handler, a config watcher, or an admin action. `RegisterService` adapts the `Reloadable` interface onto the callback, and `ActionServiceReload` broadcasts on success.

```go
r := c.ServiceReload(c.Context())
if !r.OK {
	// first failing reload stops the chain
}
```
