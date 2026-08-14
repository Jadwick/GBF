//go:build windows
package main

import "fmt"
import "os"
import "os/exec"

func doesTerminalSupportColor() bool {
	// Check if a powershell variable is set in Env
    s, b := os.LookupEnv("PsModulePath")
	fmt.Println(s)
    return b
}

func systemClearscreen() {
	exec.Command("cmd", "/C", "cls")
}
