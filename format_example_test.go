// SPDX-License-Identifier: EUPL-1.2

package core_test

import . "dappco.re/go"

func ExampleSprintln() {
	Print(Stdout(), Sprintln("deploy", "done"))
	// Output: deploy done
}

func ExampleErrorf() {
	err := Errorf("deploy failed: %s", "timeout")
	Println(err.Error())
	// Output: deploy failed: timeout
}

// ExampleSprint concatenates operands into a string through `Sprint`.
func ExampleSprint() {
	Println(Sprint("count=", 42))
	// Output: count=42
}

// ExampleSprintf formats a string through `Sprintf`.
func ExampleSprintf() {
	Println(Sprintf("%s=%d", "agents", 5))
	// Output: agents=5
}

// ExamplePrintln writes a line to stdout through `Println`.
func ExamplePrintln() {
	Println("agent ready")
	// Output: agent ready
}

// ExamplePrint writes formatted output to a writer through `Print`.
func ExamplePrint() {
	buf := NewBuffer()
	Print(buf, "%s ready", "agent")
	Println(buf.String())
	// Output: agent ready
}
