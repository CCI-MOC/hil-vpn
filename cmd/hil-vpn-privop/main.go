package main

import (
	"fmt"
	"os"
	"strings"
)

// Print a help message to stderr and exit with the given status code.
func usage(exitCode int) {
	fmt.Fprintln(os.Stderr, strings.Join([]string{
		`Usage:`,
		``,
		`    hil-vpn-privop create <name> <vlan-no> <port-no>`,
		`    hil-vpn-privop start <name>`,
		`    hil-vpn-privop stop <name>`,
		`    hil-vpn-privop delete <name>`,
		`    hil-vpn-privop list`,
	}, "\n",
	))
	os.Exit(exitCode)
}

// Verify that the number of subcommand-specific arguments is equal to count.
// If not, prints a help message and exits with a failing status code.
func checkNumArgs(count int) {
	if len(os.Args) != count+2 {
		fmt.Fprintf(os.Stderr, "Wrong number of arguments for subcommand %q\n\n", os.Args[1])
		usage(1)
	}
}

func main() {
	if len(os.Args) < 2 {
		usage(1)
	}
	switch os.Args[1] {
	case "create":
		checkNumArgs(3)
		fmt.Fprintln(os.Stderr, "Unimplemented")
	case "start":
		checkNumArgs(1)
		fmt.Fprintln(os.Stderr, "Unimplemented")
	case "stop":
		checkNumArgs(1)
		fmt.Fprintln(os.Stderr, "Unimplemented")
	case "delete":
		checkNumArgs(1)
		fmt.Fprintln(os.Stderr, "Unimplemented")
	case "list":
		checkNumArgs(0)
		fmt.Fprintln(os.Stderr, "Unimplemented")
	case "-h", "--help", "help":
		usage(0)
	default:
		fmt.Fprintf(os.Stderr, "Unknown subcommand: %q\n\n", os.Args[1])
		usage(1)
	}
}
