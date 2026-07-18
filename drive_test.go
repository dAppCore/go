package core_test

import (
	. "dappco.re/go"
)

// --- Drive (Transport Handles) ---

func TestDrive_New_Good(t *T) {
	c := New()
	r := c.Drive().New(NewOptions(
		Option{Key: "name", Value: "api"},
		Option{Key: "transport", Value: "https://api.lthn.ai"},
	))
	AssertTrue(t, r.OK)
	AssertEqual(t, "api", r.Value.(*DriveHandle).Name)
	AssertEqual(t, "https://api.lthn.ai", r.Value.(*DriveHandle).Transport)
}

func TestDrive_New_Bad(t *T) {
	c := New()
	// Missing name
	r := c.Drive().New(NewOptions(
		Option{Key: "transport", Value: "https://api.lthn.ai"},
	))
	AssertFalse(t, r.OK)
}

func TestDrive_Get_Good(t *T) {
	c := New()
	c.Drive().New(NewOptions(
		Option{Key: "name", Value: "ssh"},
		Option{Key: "transport", Value: "ssh://claude@10.69.69.165"},
	))
	r := c.Drive().Get("ssh")
	AssertTrue(t, r.OK)
	handle := r.Value.(*DriveHandle)
	AssertEqual(t, "ssh://claude@10.69.69.165", handle.Transport)
}

func TestDrive_Get_Bad(t *T) {
	c := New()
	r := c.Drive().Get("nonexistent")
	AssertFalse(t, r.OK)
}

func TestDrive_Has_Good(t *T) {
	c := New()
	c.Drive().New(NewOptions(Option{Key: "name", Value: "mcp"}, Option{Key: "transport", Value: "mcp://mcp.lthn.sh"}))
	AssertTrue(t, c.Drive().Has("mcp"))
	AssertFalse(t, c.Drive().Has("missing"))
}

func TestDrive_Names_Good(t *T) {
	c := New()
	c.Drive().New(NewOptions(Option{Key: "name", Value: "api"}, Option{Key: "transport", Value: "https://api.lthn.ai"}))
	c.Drive().New(NewOptions(Option{Key: "name", Value: "ssh"}, Option{Key: "transport", Value: "ssh://claude@10.69.69.165"}))
	c.Drive().New(NewOptions(Option{Key: "name", Value: "mcp"}, Option{Key: "transport", Value: "mcp://mcp.lthn.sh"}))
	names := c.Drive().Names()
	AssertLen(t, names, 3)
	AssertContains(t, names, "api")
	AssertContains(t, names, "ssh")
	AssertContains(t, names, "mcp")
}

func TestDrive_Drive_New_OptionsPreserved_Good(t *T) {
	c := New()
	c.Drive().New(NewOptions(
		Option{Key: "name", Value: "api"},
		Option{Key: "transport", Value: "https://api.lthn.ai"},
		Option{Key: "timeout", Value: 30},
	))
	r := c.Drive().Get("api")
	AssertTrue(t, r.OK)
	handle := r.Value.(*DriveHandle)
	AssertEqual(t, 30, handle.Options.Int("timeout"))
}

// --- AX-7 canonical triplets ---

func TestDrive_Drive_New_Good(t *T) {
	c := New()
	r := c.Drive().New(NewOptions(
		Option{Key: "name", Value: "homelab"},
		Option{Key: "transport", Value: "ssh://agent@10.69.69.165"},
	))
	AssertTrue(t, r.OK)
	handle := r.Value.(*DriveHandle)
	AssertEqual(t, "homelab", handle.Name)
	AssertEqual(t, "ssh://agent@10.69.69.165", handle.Transport)
}

func TestDrive_Drive_New_Bad(t *T) {
	c := New()
	r := c.Drive().New(NewOptions(Option{Key: "transport", Value: "mcp://agent.local"}))
	AssertFalse(t, r.OK)
}

func TestDrive_Drive_New_Ugly(t *T) {
	c := New()
	r := c.Drive().New(NewOptions(Option{Key: "name", Value: "loopback"}))
	AssertTrue(t, r.OK)
	handle := r.Value.(*DriveHandle)
	AssertEqual(t, "", handle.Transport)
	AssertTrue(t, c.Drive().Has("loopback"))
}

// --- Named handles: c.Drive(name) / On / Exists / Handle / Transport ---

func TestDrive_On_Good(t *T) {
	c := New()
	c.Drive().New(NewOptions(
		Option{Key: "name", Value: "homelab"},
		Option{Key: "transport", Value: "ssh://agent@10.69.69.165"},
	))
	AssertTrue(t, c.Drive("homelab").Exists())
}

func TestDrive_On_Bad(t *T) {
	AssertFalse(t, New().Drive("ghost").Exists())
}

func TestDrive_On_Ugly(t *T) {
	c := New()
	// An empty name returns the root subsystem, same as zero-arg.
	AssertSame(t, c.Drive(), c.Drive(""))
	// The unbound subsystem is never a handle.
	AssertFalse(t, c.Drive().Exists())
}

func TestDrive_Handle_Good(t *T) {
	c := New()
	c.Drive().New(NewOptions(
		Option{Key: "name", Value: "forge"},
		Option{Key: "transport", Value: "https://api.lthn.ai"},
	))
	r := c.Drive("forge").Handle()
	AssertTrue(t, r.OK)
	AssertEqual(t, "forge", r.Value.(*DriveHandle).Name)
}

func TestDrive_Handle_Bad(t *T) {
	r := New().Drive().Handle()
	AssertFalse(t, r.OK)
	AssertContains(t, r.Error(), "no handle bound")
}

func TestDrive_Handle_Ugly(t *T) {
	r := New().Drive("ghost").Handle()
	AssertFalse(t, r.OK)
	AssertContains(t, r.Error(), "not found")
}

func TestDrive_Transport_Good(t *T) {
	c := New()
	c.Drive().New(NewOptions(
		Option{Key: "name", Value: "homelab"},
		Option{Key: "transport", Value: "ssh://agent@10.69.69.165"},
	))
	AssertEqual(t, "ssh://agent@10.69.69.165", c.Drive("homelab").Transport())
}

func TestDrive_Transport_Bad(t *T) {
	AssertEqual(t, "", New().Drive("ghost").Transport())
}

func TestDrive_Transport_Ugly(t *T) {
	AssertEqual(t, "", New().Drive().Transport())
}

func TestDrive_Drive_Exists_Good(t *T) {
	c := New()
	c.Drive().New(NewOptions(
		Option{Key: "name", Value: "homelab"},
		Option{Key: "transport", Value: "ssh://agent@10.69.69.165"},
	))
	AssertTrue(t, c.Drive("homelab").Exists())
}

func TestDrive_Drive_Exists_Bad(t *T) {
	AssertFalse(t, New().Drive("ghost").Exists())
}

func TestDrive_Drive_Exists_Ugly(t *T) {
	// The unbound subsystem never answers as a handle.
	AssertFalse(t, New().Drive().Exists())
}
