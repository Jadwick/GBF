//go:build linux
package main

import "os/exec"

func systemClearscreen() {
	exec.Command("clear")
}
