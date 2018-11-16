package main

import (
	"os/exec"
)

// The actually functionality of each of the commands; the function
// <command>Cmd implements the named command. The caller is responsible
// for validating hte argumetns. If an an error occurs, the program will
// exit with a failing status code.

// Implement the 'create' subcommand.
func createCmd(vpnName string, vlanNo, portNo uint16) string {
	cfg, err := NewOpenVpnConfig(vpnName, portNo)
	chkfatal("Generating openvpn config:", err)
	chkfatal("Saving openvpn config:", cfg.Save())
	return cfg.Key
}

// Implement the 'start' subcommand.
func startCmd(vpnName string) {
	err := exec.Command("systemctl", "enable", "--now", "openvpn-server@"+vpnName).Run()
	chkfatal("Starting & enabling vpn", err)
}

// Implement the 'stop' subcommand.
func stopCmd(vpnName string) {
	err := exec.Command("systemctl", "disable", "--now", "openvpn-server@"+vpnName).Run()
	chkfatal("Stopping & disabling vpn", err)
}
