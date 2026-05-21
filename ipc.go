// SPDX-License-Identifier: EUPL-1.2

// Message bus for the Core framework.
// Dispatches actions (fire-and-forget), queries (first responder),
// and tasks (first executor) between registered handlers.

package core

// actionHandlers + queryHandlers are slice types stored behind
// AtomicPointer. Dispatch (broadcast / Query / QueryAll) loads the
// pointer lock-free and iterates the slice directly — no per-call
// clone, no RLock contention. Registration takes a short mutex while
// it copies-on-write to a new slice and stores the pointer back.
type actionHandlers []func(*Core, Message) Result
type queryHandlers []QueryHandler

// Ipc holds IPC dispatch data and the named action registry.
//
//	ipc := (&core.Ipc{}).New()
type Ipc struct {
	ipcRegMu   RWMutex                       // serialises action registration
	ipcActions AtomicPointer[actionHandlers] // dispatch reads this lock-free

	queryRegMu RWMutex                      // serialises query registration
	queryFns   AtomicPointer[queryHandlers] // dispatch reads this lock-free

	actions *Registry[*Action] // named action registry
	tasks   *Registry[*Task]   // named task registry
}

// broadcast dispatches a message to all registered IPC handlers.
// Each handler is wrapped in panic recovery. All handlers fire regardless of individual results.
func (c *Core) broadcast(msg Message) Result {
	handlers := c.ipc.ipcActions.Load()
	if handlers == nil {
		return Result{OK: true}
	}
	for _, h := range *handlers {
		func() {
			defer func() {
				if r := recover(); r != nil {
					Error("ACTION handler panicked", "panic", r)
				}
			}()
			h(c, msg)
		}()
	}
	return Result{OK: true}
}

// Query dispatches a request — first handler to return OK wins.
//
//	r := c.Query(MyQuery{})
func (c *Core) Query(q Query) Result {
	handlers := c.ipc.queryFns.Load()
	if handlers == nil {
		return Result{}
	}
	for _, h := range *handlers {
		r := h(c, q)
		if r.OK {
			return r
		}
	}
	return Result{}
}

// QueryAll dispatches a request — collects all OK responses.
//
//	r := c.QueryAll(countQuery{})
//	results := r.Value.([]any)
func (c *Core) QueryAll(q Query) Result {
	handlers := c.ipc.queryFns.Load()
	if handlers == nil {
		return Result{[]any(nil), true}
	}
	var results []any
	for _, h := range *handlers {
		r := h(c, q)
		if r.OK && r.Value != nil {
			results = append(results, r.Value)
		}
	}
	return Result{results, true}
}

// RegisterQuery registers a handler for QUERY dispatch.
//
//	c.RegisterQuery(func(_ *core.Core, q core.Query) core.Result { ... })
func (c *Core) RegisterQuery(handler QueryHandler) {
	c.ipc.queryRegMu.Lock()
	defer c.ipc.queryRegMu.Unlock()
	cur := c.ipc.queryFns.Load()
	var next queryHandlers
	if cur != nil {
		next = make(queryHandlers, len(*cur)+1)
		copy(next, *cur)
		next[len(*cur)] = handler
	} else {
		next = queryHandlers{handler}
	}
	c.ipc.queryFns.Store(&next)
}

// --- IPC Registration (handlers) ---

// RegisterAction registers a broadcast handler for ACTION messages.
//
//	c.RegisterAction(func(c *core.Core, msg core.Message) core.Result {
//	    if ev, ok := msg.(AgentCompleted); ok { ... }
//	    return core.Result{OK: true}
//	})
func (c *Core) RegisterAction(handler func(*Core, Message) Result) {
	c.ipc.ipcRegMu.Lock()
	defer c.ipc.ipcRegMu.Unlock()
	cur := c.ipc.ipcActions.Load()
	var next actionHandlers
	if cur != nil {
		next = make(actionHandlers, len(*cur)+1)
		copy(next, *cur)
		next[len(*cur)] = handler
	} else {
		next = actionHandlers{handler}
	}
	c.ipc.ipcActions.Store(&next)
}

// RegisterActions registers multiple broadcast handlers.
//
//	c := core.New()
//	c.RegisterActions(
//	    func(c *core.Core, msg core.Message) core.Result { return core.Result{OK: true} },
//	    func(c *core.Core, msg core.Message) core.Result { return core.Result{OK: true} },
//	)
func (c *Core) RegisterActions(handlers ...func(*Core, Message) Result) {
	if len(handlers) == 0 {
		return
	}
	c.ipc.ipcRegMu.Lock()
	defer c.ipc.ipcRegMu.Unlock()
	cur := c.ipc.ipcActions.Load()
	var next actionHandlers
	if cur != nil {
		next = make(actionHandlers, len(*cur)+len(handlers))
		copy(next, *cur)
		copy(next[len(*cur):], handlers)
	} else {
		next = make(actionHandlers, len(handlers))
		copy(next, handlers)
	}
	c.ipc.ipcActions.Store(&next)
}
