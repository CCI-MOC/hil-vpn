package main

import (
	"fmt"
	"os"
	"os/exec"
)

// The actual functionality of the 'create' subcommand. The caller is
// responsible for validating the arguments. If an error occurs,
// the program will exit with an error. Otherwise, the generated
// static key is returned.
func createCmd(vpnName string, vlanNo, portNo uint16) string {
	cmd := exec.Command("openvpn", "--genkey", "--secret", "/dev/fd/1")
	output, err := cmd.Output()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error invoking openvpn:", err)
		os.Exit(1)
	}
	return string(output)
}
