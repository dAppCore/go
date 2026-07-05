package core_test

import (
	. "dappco.re/go"
)

// --- IPC: named action invocation ---

func TestIpc_Core_Action_Good(t *T) {
	c := New()
	c.Action("compute", func(_ Context, opts Options) Result {
		return Result{Value: 42, OK: true}
	})
	r := c.Action("compute").Run(Background(), NewOptions())
	AssertTrue(t, r.OK)
	AssertEqual(t, 42, r.Value)
}

// --- IPC: queries ---

func TestIpc_Core_Query_Good(t *T) {
	c := New()
	c.RegisterQuery(func(_ *Core, q Query) Result {
		if q == "agent.status" {
			return Result{Value: "ready", OK: true}
		}
		return Result{}
	})

	r := c.Query("agent.status")

	AssertTrue(t, r.OK)
	AssertEqual(t, "ready", r.Value)
}

func TestIpc_Core_Query_Bad(t *T) {
	r := New().Query("missing")
	AssertFalse(t, r.OK)
}

func TestIpc_Core_Query_Ugly(t *T) {
	c := New()
	c.RegisterQuery(func(_ *Core, _ Query) Result {
		return Result{Value: nil, OK: true}
	})

	r := c.Query(nil)

	AssertTrue(t, r.OK)
	AssertNil(t, r.Value)
}

func TestIpc_Core_QueryAll_Good(t *T) {
	c := New()
	c.RegisterQuery(func(_ *Core, _ Query) Result { return Result{Value: "agent", OK: true} })
	c.RegisterQuery(func(_ *Core, _ Query) Result { return Result{Value: "health", OK: true} })

	r := c.QueryAll("status")

	AssertTrue(t, r.OK)
	AssertElementsMatch(t, []any{"agent", "health"}, r.Value.([]any))
}

func TestIpc_Core_QueryAll_Bad(t *T) {
	c := New()
	c.RegisterQuery(func(_ *Core, _ Query) Result { return Result{Value: "ignored", OK: false} })

	r := c.QueryAll("status")

	AssertTrue(t, r.OK)
	AssertEmpty(t, r.Value.([]any))
}

func TestIpc_Core_QueryAll_Ugly(t *T) {
	c := New()
	c.RegisterQuery(func(_ *Core, _ Query) Result { return Result{Value: nil, OK: true} })

	r := c.QueryAll(nil)

	AssertTrue(t, r.OK)
	AssertEmpty(t, r.Value.([]any))
}

// --- IPC: actions (register + broadcast dispatch) ---

func TestIpc_Core_RegisterAction_Good(t *T) {
	c := New()
	var received Message
	c.RegisterAction(func(_ *Core, msg Message) Result {
		received = msg
		return Result{OK: true}
	})

	AssertTrue(t, c.ACTION("agent.dispatch").OK)
	AssertEqual(t, Message("agent.dispatch"), received)
}

func TestIpc_Core_RegisterAction_Bad(t *T) {
	c := New()
	c.RegisterAction(nil)
	c.RegisterAction(func(_ *Core, _ Message) Result {
		return Result{Value: NewError("handler fails"), OK: false}
	})

	// Broadcast ignores handler results: a nil or failing handler must not
	// fail the dispatch.
	r := c.ACTION("agent.dispatch")

	AssertTrue(t, r.OK)
}

func TestIpc_Core_RegisterAction_Ugly(t *T) {
	c := New()
	var order []int
	c.RegisterAction(func(_ *Core, _ Message) Result {
		order = append(order, 1)
		return Result{OK: true}
	})
	c.RegisterAction(func(_ *Core, _ Message) Result {
		panic("handler failed")
	})
	c.RegisterAction(func(_ *Core, _ Message) Result {
		order = append(order, 3)
		return Result{OK: true}
	})

	AssertTrue(t, c.ACTION(nil).OK)
	AssertEqual(t, []int{1, 3}, order, "a panicking handler must not stop the chain")
}

func TestIpc_Core_RegisterActions_Good(t *T) {
	c := New()
	count := 0
	handler := func(_ *Core, _ Message) Result { count++; return Result{OK: true} }

	c.RegisterActions(handler, handler)
	AssertTrue(t, c.ACTION("agent.dispatch").OK)

	AssertEqual(t, 2, count)
}

func TestIpc_Core_RegisterActions_Bad(t *T) {
	c := New()
	c.RegisterActions()

	r := c.ACTION("agent.dispatch")

	AssertTrue(t, r.OK)
}

func TestIpc_Core_RegisterActions_Ugly(t *T) {
	c := New()
	count := 0
	c.RegisterActions(
		nil,
		func(_ *Core, _ Message) Result { count++; return Result{OK: true} },
	)

	r := c.ACTION(nil)

	AssertTrue(t, r.OK)
	AssertEqual(t, 1, count)
}

func TestIpc_Core_RegisterQuery_Good(t *T) {
	c := New()
	c.RegisterQuery(func(_ *Core, _ Query) Result {
		return Result{Value: "ready", OK: true}
	})

	r := c.Query("agent.status")

	AssertTrue(t, r.OK)
	AssertEqual(t, "ready", r.Value)
}

func TestIpc_Core_RegisterQuery_Bad(t *T) {
	c := New()
	c.RegisterQuery(func(_ *Core, _ Query) Result {
		return Result{Value: "ignored", OK: false}
	})

	r := c.Query("agent.status")

	AssertFalse(t, r.OK)
}

func TestIpc_Core_RegisterQuery_Ugly(t *T) {
	c := New()
	c.RegisterQuery(func(_ *Core, _ Query) Result { return Result{} })
	c.RegisterQuery(func(_ *Core, _ Query) Result { return Result{Value: "fallback", OK: true} })

	r := c.Query(nil)

	AssertTrue(t, r.OK)
	AssertEqual(t, "fallback", r.Value)
}
