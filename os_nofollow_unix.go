// SPDX-License-Identifier: EUPL-1.2

//go:build !windows

// O_NOFOLLOW open flag — unix value. Split into a build-tagged file
// because the constant lives in syscall (not os) and Windows has no
// equivalent. syscall is banned everywhere in core/go EXCEPT this single
// declaration (per #1681) — consumers can't reach syscall.O_NOFOLLOW
// themselves, so core re-exports it here.

package core

import "syscall"

// O_NOFOLLOW makes OpenFile fail with a symlink error if the final path
// component is a symbolic link — refusing to open a symlink-resolving
// path. Combine with the other O_* flags. On Windows this constant is 0
// (no-op); the symlink-refusal guarantee is unix-only.
//
//	// Refuse to create-or-open through a symlinked final component.
//	r := core.OpenFile(p, core.O_CREATE|core.O_EXCL|core.O_NOFOLLOW|core.O_WRONLY, 0o600)
//	if !r.OK { return r }
const O_NOFOLLOW = syscall.O_NOFOLLOW
