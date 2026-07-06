// SPDX-License-Identifier: EUPL-1.2

// Synchronisation, locking, and lifecycle snapshots for the Core framework.

package core

// Lock is the DTO for a named mutex — pure data, no embedded cache.
// The per-Core cache of *Lock wrappers lives in Core's locks Registry.
//
// Mutex is the backing core.RWMutex.
//
//	c := core.New()
//	lock := c.Lock("service-registry")
//	lock.Lock(); defer lock.Unlock()
type Lock struct {
	Name  string
	Mutex *RWMutex
}

// Lock returns a named Lock, creating the mutex if needed.
// Locks are per-Core — separate Core instances do not share mutexes.
// The returned *Lock wrapper is cached in Core's locks Registry, so
// subsequent calls for the same name reuse the same pointer (race-safe via
// Registry.GetOrSet — two goroutines racing the first lookup converge on a
// single shared *RWMutex).
//
//	l := c.Lock("drain")
//	l.Lock(); defer l.Unlock()
func (c *Core) Lock(name string) *Lock {
	return c.locks.GetOrSet(name, func() *Lock {
		return &Lock{Name: name, Mutex: &RWMutex{}}
	}).Value.(*Lock)
}

// Lock acquires the named mutex for write.
//
//	c.Lock("drain").Lock()
//	defer c.Lock("drain").Unlock()
func (l *Lock) Lock() { l.Mutex.Lock() }

// Unlock releases the named mutex from write.
//
//	c.Lock("drain").Unlock()
func (l *Lock) Unlock() { l.Mutex.Unlock() }

// RLock acquires the named mutex for read.
//
//	c.Lock("cache").RLock()
//	defer c.Lock("cache").RUnlock()
func (l *Lock) RLock() { l.Mutex.RLock() }

// RUnlock releases the named mutex from read.
//
//	c.Lock("cache").RUnlock()
func (l *Lock) RUnlock() { l.Mutex.RUnlock() }

// TryLock attempts to acquire the write mutex without blocking.
// Returns Result{OK: true} when acquired, Result{OK: false} when held.
//
//	if c.Lock("drain").TryLock().OK {
//	    defer c.Lock("drain").Unlock()
//	    // ...
//	}
func (l *Lock) TryLock() Result {
	return l.Mutex.TryLock()
}

// LockEnable marks that the service lock should be applied after
// initialisation. The lock is service-global — it freezes the whole services
// registry against further registration; there is no per-service lock (a
// service is a registry entry, not a lockable registry). Pair with LockApply.
//
//	c := core.New()
//	c.LockEnable()
//	c.LockApply()
func (c *Core) LockEnable() {
	c.services.lockEnabled = true
}

// LockApply activates the service-global lock if it was enabled via LockEnable
// (or WithServiceLock). A no-op when the lock was never enabled.
//
//	c := core.New(core.WithServiceLock())
//	c.LockApply()
func (c *Core) LockApply() {
	if c.services.lockEnabled {
		c.services.Lock()
	}
}

// Startables returns services that have an OnStart function, in registration order.
//
//	c := core.New()
//	r := c.Startables()
//	if r.OK { services := r.Value.([]*core.Service); _ = services }
func (c *Core) Startables() Result {
	if c.services == nil {
		return Result{}
	}
	out := make([]*Service, 0, c.services.Len())
	c.services.Each(func(_ string, svc *Service) {
		if svc.OnStart != nil {
			out = append(out, svc)
		}
	})
	if len(out) == 0 {
		out = nil // byte-identical to the old var-nil behaviour
	}
	return Result{out, true}
}

// Stoppables returns services that have an OnStop function, in registration order.
//
//	c := core.New()
//	r := c.Stoppables()
//	if r.OK { services := r.Value.([]*core.Service); _ = services }
func (c *Core) Stoppables() Result {
	if c.services == nil {
		return Result{}
	}
	out := make([]*Service, 0, c.services.Len())
	c.services.Each(func(_ string, svc *Service) {
		if svc.OnStop != nil {
			out = append(out, svc)
		}
	})
	if len(out) == 0 {
		out = nil // byte-identical to the old var-nil behaviour
	}
	return Result{out, true}
}
