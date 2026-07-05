// SPDX-License-Identifier: EUPL-1.2

package core

import "embed"

//go:embed all:tests/data
var EmbeddedTestFS embed.FS

func MustCompressTestAsset(t *T, input string) string {
	t.Helper()

	r := compress(input)
	RequireTrue(t, r.OK)
	return r.Value.(string)
}

func TestEmbed_Embed_path_Good(t *T) {
	embed := &Embed{basedir: "assets"}

	r := embed.path("agent/readme.md")

	AssertTrue(t, r.OK)
	AssertEqual(t, "assets/agent/readme.md", r.Value)
}
func TestEmbed_Embed_path_Bad(t *T) {
	embed := &Embed{basedir: "assets"}

	r := embed.path("../../secrets/token")

	AssertFalse(t, r.OK)
	AssertContains(t, r.Error(), "path traversal rejected")
}
func TestEmbed_Embed_path_Ugly(t *T) {
	embed := &Embed{basedir: "assets"}

	r := embed.path(".")

	AssertTrue(t, r.OK)
	AssertEqual(t, "assets", r.Value)
}
func TestEmbed_compress_Good(t *T) {
	packed := compress("agent dispatch ready")
	RequireTrue(t, packed.OK)

	plain := decompress(packed.Value.(string))

	RequireTrue(t, plain.OK)
	AssertEqual(t, "agent dispatch ready", plain.Value)
}
func TestEmbed_compress_Bad(t *T) {
	packed := compress("")
	RequireTrue(t, packed.OK)

	plain := decompress(packed.Value.(string))

	RequireTrue(t, plain.OK)
	AssertEqual(t, "", plain.Value)
}
func TestEmbed_compress_Ugly(t *T) {
	input := Join("\n", "agent", "dispatch", "retry")
	packed := compress(input)
	RequireTrue(t, packed.OK)

	plain := decompress(packed.Value.(string))

	RequireTrue(t, plain.OK)
	AssertEqual(t, input, plain.Value)
}
func TestEmbed_compressFile_Good(t *T) {
	path := Path(t.TempDir(), "agent.txt")
	RequireTrue(t, WriteFile(path, []byte("ready"), 0o644).OK)

	packed := compressFile(path)
	RequireTrue(t, packed.OK)
	plain := decompress(packed.Value.(string))

	RequireTrue(t, plain.OK)
	AssertEqual(t, "ready", plain.Value)
}
func TestEmbed_compressFile_Bad(t *T) {
	AssertFalse(t, compressFile(Path(t.TempDir(), "missing.txt")).OK)
}
func TestEmbed_compressFile_Ugly(t *T) {
	path := Path(t.TempDir(), "empty.txt")
	RequireTrue(t, WriteFile(path, nil, 0o644).OK)

	packed := compressFile(path)
	RequireTrue(t, packed.OK)
	plain := decompress(packed.Value.(string))

	RequireTrue(t, plain.OK)
	AssertEqual(t, "", plain.Value)
}
func TestEmbed_decompress_Good(t *T) {
	packed := compress("homelab")
	RequireTrue(t, packed.OK)

	plain := decompress(packed.Value.(string))

	RequireTrue(t, plain.OK)
	AssertEqual(t, "homelab", plain.Value)
}
func TestEmbed_decompress_Bad(t *T) {
	AssertFalse(t, decompress("not base64").OK)
}
func TestEmbed_decompress_Ugly(t *T) {
	AssertFalse(t, decompress(Base64Encode([]byte("plain text"))).OK)
}
func TestEmbed_getAllFiles_Good(t *T) {
	dir := t.TempDir()
	agent := Path(dir, "agent.txt")
	task := Path(dir, "nested", "task.txt")
	RequireTrue(t, WriteFile(agent, []byte("agent"), 0o644).OK)
	RequireTrue(t, MkdirAll(Path(dir, "nested"), 0o755).OK)
	RequireTrue(t, WriteFile(task, []byte("task"), 0o644).OK)

	r := getAllFiles(dir)

	RequireTrue(t, r.OK)
	files := r.Value.([]string)
	AssertContains(t, files, agent)
	AssertContains(t, files, task)
}
func TestEmbed_getAllFiles_Bad(t *T) {
	AssertFalse(t, getAllFiles(Path(t.TempDir(), "missing")).OK)
}
func TestEmbed_getAllFiles_Ugly(t *T) {
	r := getAllFiles(t.TempDir())

	RequireTrue(t, r.OK)
	AssertEmpty(t, r.Value)
}
func TestEmbed_isTemplate_Good(t *T) {
	AssertTrue(t, isTemplate("README.md.tmpl", []string{".tmpl"}))
}
func TestEmbed_isTemplate_Bad(t *T) {
	AssertFalse(t, isTemplate("README.md", []string{".tmpl"}))
}
func TestEmbed_isTemplate_Ugly(t *T) {
	AssertTrue(t, isTemplate("agent.go.tpl", []string{".tmpl", ".tpl"}))
}
func TestEmbed_renderPath_Good(t *T) {
	path := renderPath("workspace/{{.Name}}/README.md", map[string]string{"Name": "agent"})

	AssertEqual(t, "workspace/agent/README.md", path)
}
func TestEmbed_renderPath_Bad(t *T) {
	path := "workspace/{{.Name/README.md"

	AssertEqual(t, path, renderPath(path, map[string]string{"Name": "agent"}))
}
func TestEmbed_renderPath_Ugly(t *T) {
	path := "workspace/{{.Name}}/README.md"

	AssertEqual(t, path, renderPath(path, nil))
}
func TestEmbed_copyFile_Good(t *T) {
	src := t.TempDir()
	target := Path(t.TempDir(), "agent.txt")
	RequireTrue(t, WriteFile(Path(src, "agent.txt"), []byte("ready"), 0o644).OK)

	r := copyFile(DirFS(src), "agent.txt", target)

	RequireTrue(t, r.OK)
	AssertEqual(t, "ready", string(ReadFile(target).Value.([]byte)))
}
func TestEmbed_copyFile_Bad(t *T) {
	AssertFalse(t, copyFile(DirFS(t.TempDir()), "missing.txt", Path(t.TempDir(), "out.txt")).OK)
}
func TestEmbed_copyFile_Ugly(t *T) {
	src := t.TempDir()
	target := Path(t.TempDir(), "nested", "agent.txt")
	RequireTrue(t, WriteFile(Path(src, "agent.txt"), []byte("nested"), 0o644).OK)

	r := copyFile(DirFS(src), "agent.txt", target)

	RequireTrue(t, r.OK)
	AssertEqual(t, "nested", string(ReadFile(target).Value.([]byte)))
}
