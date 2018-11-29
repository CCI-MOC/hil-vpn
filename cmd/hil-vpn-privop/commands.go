package main

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
)

// The actual functionality of each of the commands; the function
// <command>Cmd implements the named command. The caller is responsible
// for validating the arguments. If an an error occurs, the program will
// exit with a failing status code.

var keyFileRe = regexp.MustCompile("^hil-vpn-([-_a-zA-Z0-9]+).key$")

// Implement the 'create' subcommand.
func createCmd(vpnName string, vlanNo, portNo uint16) string {
	cfg, err := NewOpenVpnConfig(vpnName, portNo)
	chkfatal("Generating openvpn config:", err)
	chkfatal("Saving openvpn config:", cfg.Save())
	return cfg.Key
}

// Implement the 'start' subcommand.
func startCmd(vpnName string) {
	err := exec.Command("systemctl", "enable", "--now", getServiceName(vpnName)).Run()
	chkfatal("Starting & enabling vpn", err)
}

// Implement the 'stop' subcommand.
func stopCmd(vpnName string) {
	err := exec.Command("systemctl", "disable", "--now", getServiceName(vpnName)).Run()
	chkfatal("Stopping & disabling vpn", err)
}

// Implement the 'delete' subcommand.
func deleteCmd(vpnName string) {
	err := exec.Command("systemctl", "status", getServiceName(vpnName)).Run()
	if err == nil {
		fmt.Fprintf(
			os.Stderr,
			"Error: cannot delete vpn: %v; it is still running.",
			vpnName,
		)
		os.Exit(1)
	}
	if _, ok := err.(*exec.ExitError); !ok {
		chkfatal("Checking vpn status:", err)
	}

	// A failing exit status indicates that the service was not running; go
	// ahead and delete the config & key:

	chkfatal("Deleting vpn key file", os.Remove(getKeyPath(vpnName)))
	chkfatal("Deleting vpn config file", os.Remove(getCfgPath(vpnName)))
}

// Implement the 'list' subcommand.
func listCmd() {
	f, err := os.Open(configDir)
	chkfatal("Opening openvpn config directory", err)
	defer f.Close()
	fis, err := f.Readdir(0)
	chkfatal("Scanning openvpn config directory", err)
	for _, fi := range fis {
		matches := keyFileRe.FindStringSubmatch(fi.Name())
		if matches == nil {
			continue
		} else if len(matches) != 2 {
			panic("BUG: keyFileRe should always return a slice " +
				"of length 2: the full match and the submatch.")
		}
		fmt.Println(matches[1])
	}
}
