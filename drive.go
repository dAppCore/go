// SPDX-License-Identifier: EUPL-1.2

// Drive is the resource handle registry for transport connections.
// Packages register their transport handles (API, MCP, SSH, VPN)
// and other packages access them by name.
//
// Register a transport:
//
//	c.Drive().New(core.NewOptions(
//	    core.Option{Key: "name", Value: "api"},
//	    core.Option{Key: "transport", Value: "https://api.lthn.ai"},
//	))
//	c.Drive().New(core.NewOptions(
//	    core.Option{Key: "name", Value: "ssh"},
//	    core.Option{Key: "transport", Value: "ssh://claude@10.69.69.165"},
//	))
//	c.Drive().New(core.NewOptions(
//	    core.Option{Key: "name", Value: "mcp"},
//	    core.Option{Key: "transport", Value: "mcp://mcp.lthn.sh"},
//	))
//
// Retrieve a handle:
//
//	api := c.Drive().Get("api")
package core

// DriveHandle holds a named transport resource.
//
//	handle := &core.DriveHandle{Name: "homelab", Transport: "ssh://agent@10.69.69.165"}
//	core.Println(handle.Transport)
type DriveHandle struct {
	Name      string
	Transport string
	Options   Options
}

// Drive manages named transport handles. Embeds Registry[*DriveHandle].
//
//	c := core.New()
//	c.Drive().New(core.NewOptions(
//	    core.Option{Key: "name", Value: "homelab"},
//	    core.Option{Key: "transport", Value: "ssh://agent@10.69.69.165"},
//	))
//
// A view bound to one handle answers for it directly (see On / c.Drive(name)):
//
//	if c.Drive("homelab").Exists() { transport := c.Drive("homelab").Transport() }
type Drive struct {
	*Registry[*DriveHandle]
	bound string // non-empty on a handle-bound view — see On / c.Drive(name)
}

// On returns a view of Drive bound to a named handle. The view shares
// the handle registry with the parent; only the binding is new. Sugar
// reached via c.Drive(name).
//
//	forge := c.Drive().On("forge")
//	if forge.Exists() { ... }
func (d *Drive) On(name string) *Drive {
	return &Drive{Registry: d.Registry, bound: name}
}

// Exists reports whether the bound handle is registered — the
// capability check for named transports.
//
//	if c.Drive("forge").Exists() { ... }
func (d *Drive) Exists() bool {
	return d.bound != "" && d.Has(d.bound)
}

// Handle returns the bound transport handle. Fails with operation
// "drive.Handle" when the view has no binding or the handle is missing.
//
//	r := c.Drive("forge").Handle()
//	if r.OK { handle := r.Value.(*core.DriveHandle) }
func (d *Drive) Handle() Result {
	if d.bound == "" {
		return Result{E("drive.Handle", "no handle bound — use c.Drive(name)", nil), false}
	}
	r := d.Get(d.bound)
	if !r.OK {
		return Result{E("drive.Handle", Concat("handle not found: ", d.bound), nil), false}
	}
	return r
}

// Transport returns the bound handle's transport URL, "" when the
// binding is absent — the Options-getter contract for named transports.
//
//	url := c.Drive("homelab").Transport()  // "ssh://agent@10.69.69.165"
func (d *Drive) Transport() string {
	r := d.Handle()
	if !r.OK {
		return ""
	}
	return r.Value.(*DriveHandle).Transport
}

// New registers a transport handle.
//
//	c.Drive().New(core.NewOptions(
//	    core.Option{Key: "name", Value: "api"},
//	    core.Option{Key: "transport", Value: "https://api.lthn.ai"},
//	))
func (d *Drive) New(opts Options) Result {
	name := opts.String("name")
	if name == "" {
		return Result{}
	}

	handle := &DriveHandle{
		Name:      name,
		Transport: opts.String("transport"),
		Options:   opts,
	}

	d.Set(name, handle)
	return Result{handle, true}
}
