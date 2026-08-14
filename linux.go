//go:build linux
package main

import "os/exec"

func doesTerminalSupportColor() bool {
    //sure they do???
	return true
}

func systemClearscreen() {
	exec.Command("bash", "-c", "clear")
}
