// SPDX-License-Identifier: EUPL-1.2

package core

// Good — the user's home is protected, whichever way the OS names it.
func TestFs_isProtectedPath_Good(t *T) {
	home := UserHomeDir()
	RequireTrue(t, home.OK)

	AssertTrue(t, isProtectedPath(home.Value.(string)))
}

// Bad — an ordinary folder, even one inside the home, is not protected.
func TestFs_isProtectedPath_Bad(t *T) {
	AssertFalse(t, isProtectedPath(t.TempDir()))
	AssertFalse(t, isProtectedPath(Path(UserHomeDir().Value.(string), "logs")))
	AssertFalse(t, isProtectedPath(""))
}

// Ugly — roots are protected (a drive root on Windows), and on Windows the
// home matches regardless of case, as its filesystem does.
func TestFs_isProtectedPath_Ugly(t *T) {
	root := PathDir(UserHomeDir().Value.(string))
	for PathDir(root) != root {
		root = PathDir(root)
	}
	AssertTrue(t, isProtectedPath(root))
	if OS() == "windows" {
		AssertTrue(t, isProtectedPath(Upper(UserHomeDir().Value.(string))))
	}
}
