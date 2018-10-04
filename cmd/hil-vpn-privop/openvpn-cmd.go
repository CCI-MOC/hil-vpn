package main

import (
	"os/exec"
)

// This file provides wrappers for invoking the 'openvpn' command line tool.

// Return the path to the openvpn executable. Prints an error and
// quits if the executable cannot be found in $PATH, or if the
// command is not located at one of the paths in 'acceptablePaths'.
func findOpenVpn() string {
	return safeFindCmd("openvpn")
}

// Return a command which will execute `openvpn` with the given arguments.
// Calls 'findOpenVpn' to find the executable.
func openVpn(arg ...string) *exec.Cmd {
	path := findOpenVpn()
	return exec.Command(path, arg...)
}
