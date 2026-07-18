// SPDX-License-Identifier: EUPL-1.2

package core_test

import . "dappco.re/go"

func ExampleUserCurrent() {
	// Succeeds on a normal host; in a bare container without /etc/passwd it
	// returns OK=false. Shown without asserted output because the username and
	// home dir are host-specific.
	r := UserCurrent()
	if r.OK {
		_ = r.Value.(*User).HomeDir
	}
}

func ExampleUserLookup() {
	r := UserLookup("definitely-not-a-user-9z9z9z")
	Println(r.OK)
	// Output: false
}

func ExampleUserLookupID() {
	// Output omitted: uid existence is OS-synthesized (e.g. macOS maps high
	// uids to "nobody"), so there is no portably-absent uid to assert on.
	r := UserLookupID("0")
	if r.OK {
		_ = r.Value.(*User).Username
	}
}

func ExampleUserGroupLookup() {
	r := UserGroupLookup("definitely-not-a-group-9z9z9z")
	Println(r.OK)
	// Output: false
}
