package core_test

import . "dappco.re/go"

var testFS EmbedFS = EmbeddedTestFS

// --- Data (Embedded Content Mounts) ---

func mountTestData(t *T, c *Core, name string) {
	t.Helper()

	r := c.Data().New(NewOptions(
		Option{Key: "name", Value: name},
		Option{Key: "source", Value: testFS},
		Option{Key: "path", Value: "tests/data"},
	))
	AssertTrue(t, r.OK)
}

func TestData_New_Good(t *T) {
	c := New()
	r := c.Data().New(NewOptions(
		Option{Key: "name", Value: "test"},
		Option{Key: "source", Value: testFS},
		Option{Key: "path", Value: "tests/data"},
	))
	AssertTrue(t, r.OK)
	AssertNotNil(t, r.Value)
}

func TestData_New_Bad(t *T) {
	c := New()

	r := c.Data().New(NewOptions(Option{Key: "source", Value: testFS}))
	AssertFalse(t, r.OK)

	r = c.Data().New(NewOptions(Option{Key: "name", Value: "test"}))
	AssertFalse(t, r.OK)

	r = c.Data().New(NewOptions(Option{Key: "name", Value: "test"}, Option{Key: "source", Value: "not-an-fs"}))
	AssertFalse(t, r.OK)
}

func TestData_MountDir_Good(t *T) {
	c := New()
	fs := (&Fs{}).New("/")
	dir := MustCast[string](fs.TempDir("core-data-mountdir"))
	defer fs.DeleteAll(dir)
	fs.Write(Path(dir, "note.txt"), "hello from disk\n")

	r := c.Data().MountDir("disk", dir)
	AssertTrue(t, r.OK)

	read := c.Data("disk").ReadString("note.txt")
	AssertTrue(t, read.OK)
	AssertEqual(t, "hello from disk\n", read.Value.(string))
}

func TestData_MountDir_Bad(t *T) {
	c := New()
	r := c.Data().MountDir("", t.TempDir())
	AssertFalse(t, r.OK)

	r = c.Data().MountDir("disk", "")
	AssertFalse(t, r.OK)
}

func TestData_MountDir_Ugly(t *T) {
	c := New()
	// The parent temp dir exists; the mount target inside it does not —
	// Mount's ReadDir(".") probe fails, so MountDir fails at Mount rather
	// than succeeding structurally and merely missing on read.
	missing := Path(t.TempDir(), "does-not-exist")
	r := c.Data().MountDir("disk", missing)
	AssertFalse(t, r.OK)
}

func TestData_ReadString_Good(t *T) {
	c := New()
	mountTestData(t, c, "app")
	r := c.Data().ReadString("app/test.txt")
	AssertTrue(t, r.OK)
	AssertEqual(t, "hello from testdata\n", r.Value.(string))
}

func TestData_ReadString_Bad(t *T) {
	c := New()
	r := c.Data().ReadString("nonexistent/file.txt")
	AssertFalse(t, r.OK)
}

func TestData_ReadFile_Good(t *T) {
	c := New()
	mountTestData(t, c, "app")
	r := c.Data().ReadFile("app/test.txt")
	AssertTrue(t, r.OK)
	AssertEqual(t, "hello from testdata\n", string(r.Value.([]byte)))
}

func TestData_Get_Good(t *T) {
	c := New()
	mountTestData(t, c, "brain")
	gr := c.Data().Get("brain")
	AssertTrue(t, gr.OK)
	emb := gr.Value.(*Embed)

	r := emb.Open("test.txt")
	AssertTrue(t, r.OK)
	cr := ReadAll(r.Value)
	AssertTrue(t, cr.OK)
	AssertEqual(t, "hello from testdata\n", cr.Value)
}

func TestData_Get_Bad(t *T) {
	c := New()
	r := c.Data().Get("nonexistent")
	AssertFalse(t, r.OK)
}

func TestData_Mounts_Good(t *T) {
	c := New()
	mountTestData(t, c, "a")
	mountTestData(t, c, "b")
	mounts := c.Data().Mounts()
	AssertLen(t, mounts, 2)
}

func TestData_List_Good(t *T) {
	c := New()
	mountTestData(t, c, "app")
	r := c.Data().List("app/.")
	AssertTrue(t, r.OK)
}

func TestData_List_Bad(t *T) {
	c := New()
	r := c.Data().List("nonexistent/path")
	AssertFalse(t, r.OK)
}

func TestData_ListNames_Good(t *T) {
	c := New()
	mountTestData(t, c, "app")
	r := c.Data().ListNames("app/.")
	AssertTrue(t, r.OK)
	AssertContains(t, r.Value.([]string), "test")
}

func TestData_Extract_Good(t *T) {
	c := New()
	mountTestData(t, c, "app")
	r := c.Data().Extract("app/.", t.TempDir(), nil)
	AssertTrue(t, r.OK)
}

func TestData_Extract_Bad(t *T) {
	c := New()
	r := c.Data().Extract("nonexistent/path", t.TempDir(), nil)
	AssertFalse(t, r.OK)
}

// --- AX-7 canonical triplets ---

func TestData_Data_New_Good(t *T) {
	c := New()
	r := c.Data().New(NewOptions(
		Option{Key: "name", Value: "agent"},
		Option{Key: "source", Value: testFS},
		Option{Key: "path", Value: "tests/data"},
	))
	AssertTrue(t, r.OK)
	AssertNotNil(t, r.Value.(*Embed))
}

func TestData_Data_New_Bad(t *T) {
	c := New()
	r := c.Data().New(NewOptions(
		Option{Key: "name", Value: "agent"},
		Option{Key: "source", Value: "not-a-filesystem"},
	))
	AssertFalse(t, r.OK)
}

func TestData_Data_New_Ugly(t *T) {
	c := New()
	r := c.Data().New(NewOptions(
		Option{Key: "name", Value: "root"},
		Option{Key: "source", Value: testFS},
	))
	AssertTrue(t, r.OK)
	AssertEqual(t, ".", r.Value.(*Embed).BaseDirectory())
}

func TestData_Data_ReadFile_Good(t *T) {
	c := New()
	mountTestData(t, c, "agent")
	r := c.Data().ReadFile("agent/test.txt")
	AssertTrue(t, r.OK)
	AssertEqual(t, "hello from testdata\n", string(r.Value.([]byte)))
}

func TestData_Data_ReadFile_Bad(t *T) {
	c := New()
	r := c.Data().ReadFile("agent/test.txt")
	AssertFalse(t, r.OK)
}

func TestData_Data_ReadFile_Ugly(t *T) {
	c := New()
	mountTestData(t, c, "agent")
	r := c.Data().ReadFile("agent")
	AssertFalse(t, r.OK)
}

func TestData_Data_ReadString_Good(t *T) {
	c := New()
	mountTestData(t, c, "agent")
	r := c.Data().ReadString("agent/test.txt")
	AssertTrue(t, r.OK)
	AssertEqual(t, "hello from testdata\n", r.Value.(string))
}

func TestData_Data_ReadString_Bad(t *T) {
	c := New()
	r := c.Data().ReadString("missing/test.txt")
	AssertFalse(t, r.OK)
}

func TestData_Data_ReadString_Ugly(t *T) {
	c := New()
	mountTestData(t, c, "agent")
	r := c.Data().ReadString("")
	AssertFalse(t, r.OK)
}

func TestData_Data_List_Good(t *T) {
	c := New()
	mountTestData(t, c, "agent")
	r := c.Data().List("agent/.")
	AssertTrue(t, r.OK)
	AssertNotEmpty(t, r.Value.([]FsDirEntry))
}

func TestData_Data_List_Bad(t *T) {
	c := New()
	r := c.Data().List("missing/.")
	AssertFalse(t, r.OK)
}

func TestData_Data_List_Ugly(t *T) {
	c := New()
	mountTestData(t, c, "agent")
	r := c.Data().List("agent/test.txt")
	AssertFalse(t, r.OK)
}

func TestData_Data_ListNames_Good(t *T) {
	c := New()
	mountTestData(t, c, "agent")
	r := c.Data().ListNames("agent/.")
	AssertTrue(t, r.OK)
	AssertContains(t, r.Value.([]string), "test")
}

func TestData_Data_ListNames_Bad(t *T) {
	c := New()
	r := c.Data().ListNames("missing/.")
	AssertFalse(t, r.OK)
}

func TestData_Data_ListNames_Ugly(t *T) {
	c := New()
	mountTestData(t, c, "agent")
	r := c.Data().ListNames("agent/.")
	AssertTrue(t, r.OK)
	AssertContains(t, r.Value.([]string), "_scantest")
}

func TestData_Data_Extract_Good(t *T) {
	c := New()
	mountTestData(t, c, "agent")
	target := t.TempDir()
	r := c.Data().Extract("agent/.", target, nil)
	AssertTrue(t, r.OK)
	read := (&Fs{}).New("/").Read(Path(target, "test.txt"))
	AssertTrue(t, read.OK)
}

func TestData_Data_Extract_Bad(t *T) {
	c := New()
	r := c.Data().Extract("missing/.", t.TempDir(), nil)
	AssertFalse(t, r.OK)
}

func TestData_Data_Extract_Ugly(t *T) {
	c := New()
	mountTestData(t, c, "agent")
	target := Path(t.TempDir(), "nested", "workspace")
	r := c.Data().Extract("agent/.", target, map[string]string{"Agent": "codex"})
	AssertTrue(t, r.OK)
	AssertTrue(t, (&Fs{}).New("/").Exists(Path(target, "test.txt")).OK)
}

func TestData_Data_Mounts_Good(t *T) {
	c := New()
	mountTestData(t, c, "agent")
	mountTestData(t, c, "docs")
	AssertEqual(t, []string{"agent", "docs"}, c.Data().Mounts())
}

func TestData_Data_Mounts_Bad(t *T) {
	c := New()
	AssertEmpty(t, c.Data().Mounts())
}

func TestData_Data_Mounts_Ugly(t *T) {
	c := New()
	mountTestData(t, c, "agent")
	mounts := c.Data().Mounts()
	mounts[0] = "mutated"
	AssertEqual(t, []string{"agent"}, c.Data().Mounts())
}

// --- Named mounts: c.Data(name) / On / Exists ---

func TestData_On_Good(t *T) {
	c := New()
	mountTestData(t, c, "brain")
	r := c.Data("brain").ReadString("test.txt")
	AssertTrue(t, r.OK)
	AssertEqual(t, "hello from testdata\n", r.Value.(string))
}

func TestData_On_Bad(t *T) {
	c := New()
	// A binding to a missing mount misses on every read.
	r := c.Data("ghost").ReadString("test.txt")
	AssertFalse(t, r.OK)
}

func TestData_On_Ugly(t *T) {
	c := New()
	mountTestData(t, c, "brain")
	// The bound view and the prefixed root path read the same file.
	bound := c.Data("brain").ReadString("test.txt")
	rooted := c.Data().ReadString("brain/test.txt")
	AssertEqual(t, rooted.Value.(string), bound.Value.(string))
	// An empty name returns the root subsystem, same as zero-arg.
	AssertSame(t, c.Data(), c.Data(""))
}

func TestData_Exists_Good(t *T) {
	c := New()
	mountTestData(t, c, "brain")
	AssertTrue(t, c.Data("brain").Exists())
}

func TestData_Exists_Bad(t *T) {
	AssertFalse(t, New().Data("ghost").Exists())
}

func TestData_Exists_Ugly(t *T) {
	// The unbound subsystem is never a mount.
	AssertFalse(t, New().Data().Exists())
}
