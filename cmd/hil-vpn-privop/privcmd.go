package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// This file provides helpers for running executables, without trusting the caller.

// acceptable locations for an executable. We don't just invoke commands wherever
// they are in the $PATH, because that would allow an attacker
// who can influence the environment to run an arbitrary program with our
// permissions. Since hil-vpn-privop is intended to be used with elevated
// privileges, this is unacceptable. Instead, we whitelist common locations,
// which on a normal system are only writable by root.
var acceptablePaths = []string{
	"/bin/",
	"/sbin/",
	"/usr/bin/",
	"/usr/sbin/",
	"/usr/local/bin/",
	"/usr/local/sbin/",
}

// Return the path to an executable with the given name. Prints an error and
// quits if the executable cannot be found in $PATH, or if the command is not
// located in one of the directories in 'acceptablePaths'.
func safeFindCmd(cmdName string) string {
	path, err := exec.LookPath(cmdName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal: could not find %q executable: %v\n", cmdName, err)
		os.Exit(1)
	}
	path, err = filepath.Abs(path)
	if err != nil {
		fmt.Fprintln(
			os.Stderr,
			"Fatal: error computing path to %q executable: %v\n",
			cmdName,
			err,
		)
		os.Exit(1)
	}
	for _, v := range acceptablePaths {
		if path == v+cmdName {
			return path
		}
	}
	fmt.Fprintf(
		os.Stderr,
		"Fatal error: %q executable is in a non-standard location: %q\n",
		cmdName,
		path,
	)
	os.Exit(1)
	panic("unreachable")
}
