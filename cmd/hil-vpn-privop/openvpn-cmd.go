package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// This file provides wrappers for invoking the 'openvpn' command line tool.

// acceptable locations for an openvpn executable. We don't just invoke
// openvpn wherever it is in the $PATH, because that would allow an attacker
// who can influence the environment to run an arbitrary program with our
// permissions. Since hil-vpn-privop is intended to be used with elevated
// privileges, this is unacceptable. Instead, we whitelist common locations,
// which on a normal system are only writable by root.
var acceptablePaths = []string{
	"/bin/openvpn",
	"/sbin/openvpn",
	"/usr/bin/openvpn",
	"/usr/sbin/openvpn",
	"/usr/local/bin/openvpn",
	"/usr/local/sbin/openvpn",
}

// Return the path to the openvpn executable. Prints an error and
// quits if the executable cannot be found in $PATH, or if the
// command is not located at one of the paths in 'acceptablePaths'.
func findOpenVpn() string {
	path, err := exec.LookPath("openvpn")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Fatal: could not find openvpn executable:", err)
		os.Exit(1)
	}
	path, err = filepath.Abs(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Fatal: error finding computing path to openvpn executable:", err)
		os.Exit(1)
	}
	for _, v := range acceptablePaths {
		if path == v {
			return path
		}
	}
	fmt.Fprintln(
		os.Stderr,
		"Fatal error: openvpn executable is in a non-standard location: %q",
		path,
	)
	os.Exit(1)
	panic("unreachable")
}
