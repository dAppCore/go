// SPDX-License-Identifier: EUPL-1.2

//go:build windows

// O_NOFOLLOW open flag — Windows value. Windows has no syscall.O_NOFOLLOW
// and no symlink-refusal open semantic, so the constant is 0 (a no-op
// when OR-ed into an open flag set). The unix value lives in
// os_nofollow_unix.go. Keeping the constant defined on every platform
// lets consumers reference core.O_NOFOLLOW in portable code without a
// build tag of their own.

package core

// O_NOFOLLOW is 0 on Windows — no symlink-refusal open semantic exists.
// See os_nofollow_unix.go for the doc + usage example.
const O_NOFOLLOW = 0
