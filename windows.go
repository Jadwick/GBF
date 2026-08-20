//go:build windows
package main

import "os/exec"

func systemClearscreen() {
	exec.Command("cmd", "/c", "cls")
}
