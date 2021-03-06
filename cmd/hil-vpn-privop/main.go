package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/CCI-MOC/hil-vpn/internal/validate"
)

func init() {
	// Set the $PATH environment variable to a known-safe value. If the
	// caller is able to influence our environment, it could set PATH
	// to something containing an untrustworthy executable named openvpn
	// or systemctl. Since hil-vpn-privop runs with elevated privileges,
	// we need to guard against this, so we set PATH to a specific value,
	// rather than assuming it is sane on startup.
	err := os.Setenv(
		"PATH",
		"/bin:/sbin:/usr/bin:/usr/sbin:/usr/local/bin:/usr/local/sbin",
	)
	if err != nil {
		panic(err)
	}
}

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

// Validate that `name` is a legal name for a vpn. If so, return the name,
// otherwise exit with an error message.
func checkVpnName(name string) string {
	err := validate.CheckVpnName(name)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		usage(1)
	}
	return name
}

// Validate that `vlanStr` is a legal vlan id. If not, exit with an error
// message, otherwise parse the vlan id and return it.
func checkVlan(vlanStr string) uint16 {
	vlanNo, err := strconv.ParseInt(os.Args[3], 10, 12)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing vlan number: %v\n\n", err)
		usage(1)
	}
	if err = validate.CheckVlanNo(uint16(vlanNo)); err != nil {
		fmt.Fprintln(os.Stderr, err)
		usage(1)
	}
	return uint16(vlanNo)
}

// Validate that `portStr` is an acceptable port number. If not, exit with
// an error message, otherwise parse the port number and return it.
//
// The critera for acceptable port numbers is that they must be ports that
// a normal user could listen on; this avoids being able to use hil-vpn-privop
// to affect privileged ports.
func checkPort(portStr string) uint16 {
	portNo, err := strconv.ParseInt(os.Args[4], 10, 16)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing port number: %v\n\n", err)
		usage(1)
	}
	if portNo < 1024 {
		fmt.Fprintln(os.Stderr,
			"Unacceptable port number: ",
			portNo,
			"; only non-privileged ports (>= 1024) may be used.")
		usage(1)
	}
	return uint16(portNo)
}

func main() {
	// Make sure only one hil-vpn-privop command is running at a time:
	lockFile()
	defer unlockFile()

	if len(os.Args) < 2 {
		usage(1)
	}
	switch os.Args[1] {
	case "create":
		checkNumArgs(3)
		vpnName := checkVpnName(os.Args[2])
		vlanNo := checkVlan(os.Args[3])
		portNo := checkPort(os.Args[4])
		fmt.Print(createCmd(vpnName, vlanNo, portNo))
	case "start":
		checkNumArgs(1)
		vpnName := checkVpnName(os.Args[2])
		startCmd(vpnName)
	case "stop":
		checkNumArgs(1)
		vpnName := checkVpnName(os.Args[2])
		stopCmd(vpnName)
	case "delete":
		checkNumArgs(1)
		vpnName := checkVpnName(os.Args[2])
		deleteCmd(vpnName)
	case "list":
		checkNumArgs(0)
		listCmd()
	case "-h", "--help", "help":
		usage(0)
	default:
		fmt.Fprintf(os.Stderr, "Unknown subcommand: %q\n\n", os.Args[1])
		usage(1)
	}
}
