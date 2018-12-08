package main

import (
	"os/exec"
	"strconv"
	"strings"

	"github.com/CCI-MOC/hil-vpn/internal/staticconfig"
)

// The PrivOps interface captures the privileged operations that hil-vpnd
// needs to perform. In a production setup, we always do this by calling the
// 'hil-vpn-privop' command (see PrivOpsCmd), but this interface allows us
// to test more easily.
type PrivOps interface {
	CreateVPN(name string, vlanNo uint16, portNo uint16) (string, error)
	StartVPN(name string) error
	StopVPN(name string) error
	DeleteVPN(name string) error
	ListVPNs() ([]string, error)
}

// An implementation of PrivOps that calls the 'hil-vpn-privop' command.
type PrivOpsCmd struct{}

func privOpCmd(args ...string) *exec.Cmd {
	sudoArgs := append(
		[]string{staticconfig.Libexecdir + "/hil-vpn-privop"},
		args...,
	)
	return exec.Command("sudo", sudoArgs...)
}

func (PrivOpsCmd) CreateVPN(name string, vlanNo uint16, portNo uint16) (string, error) {
	out, err := privOpCmd(
		"create",
		name,
		strconv.Itoa(int(vlanNo)),
		strconv.Itoa(int(portNo)),
	).Output()
	return string(out), err
}

func (PrivOpsCmd) StartVPN(name string) error {
	return privOpCmd("start", name).Run()
}

func (PrivOpsCmd) StopVPN(name string) error {
	return privOpCmd("stop", name).Run()
}

func (PrivOpsCmd) DeleteVPN(name string) error {
	return privOpCmd("delete", name).Run()
}

func (PrivOpsCmd) ListVPNs() ([]string, error) {
	out, err := privOpCmd("list").Output()
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(out), "\n")

	// Chop off the trailing empty line, if any:
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}

	return lines, nil
}
