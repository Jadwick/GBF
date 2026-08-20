//go:build darwin
package main

import "os/exec"

func systemClearscreen() {
	exec.Command("clear")
}
