//go:build darwin
package main

import "os/exec"

func doesTerminalSupportColor() bool {
	return true
}

func systemClearscreen() {
	exec.Command("clear")
}
