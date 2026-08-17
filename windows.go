//go:build windows
package main

import "os"
import "os/exec"

func doesTerminalSupportColor() bool {
	// Check if a powershell variable is set in Env
    _, b := os.LookupEnv("PsModulePath")
    return b
}

func systemClearscreen() {
	exec.Command("cmd", "/c", "cls")
}
