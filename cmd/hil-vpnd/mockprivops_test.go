package main

import (
	"crypto/rand"
	"fmt"
	"sync"

	"github.com/CCI-MOC/hil-vpn/internal/validate"
)

// A mock implementation of PrivOps, for testing. Use NewMockPrivOps to
// create one of these; the zero value is not useful.
type MockPrivOps struct {
	lock sync.Mutex
	vpns map[string]*vpnInfo
}

// Create a new MockPrivOps, with no existent vpns.
func NewMockPrivOps() *MockPrivOps {
	return &MockPrivOps{
		vpns: make(map[string]*vpnInfo),
	}
}

// Info about a vpn.
type vpnInfo struct {
	// The port number that openvpn would listen on
	portNo uint16

	// The vlan number for the network
	vlanNo uint16

	// The OpenVPN static key. For testing we just use a random
	// string here.
	key string

	// Whether the vpn is up
	running bool
}

// Generate a mock "key".
func genMockKey() (string, error) {
	var key [128 / 8]byte
	_, err := rand.Read(key[:])
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", key), nil
}

// Get a vpnInfo for the named vpn, panicking if it doesn't exist.
func (ops *MockPrivOps) mustGetVpn(name string) *vpnInfo {
	vpn, ok := ops.vpns[name]
	if !ok {
		panic(fmt.Sprintf("Tried to operate on non-existent vpn %q", name))
	}
	return vpn
}

//// Implementations of the methods needed to implement the PrivOps interface.

func (ops *MockPrivOps) CreateVPN(name string, vlanNo uint16, portNo uint16) (string, error) {
	if err := validate.CheckVpnName(name); err != nil {
		panic(err)
	}
	ops.startOp()
	defer ops.endOp()
	if _, ok := ops.vpns[name]; ok {
		panic(fmt.Sprintf(
			"Tried to create a vpn with the same name (%q) as an existing one.",
			name,
		))
	}

	key, err := genMockKey()
	if err != nil {
		return "", err
	}

	ops.vpns[name] = &vpnInfo{
		portNo: portNo,
		vlanNo: vlanNo,
		key:    key,
	}

	return key, nil
}

func (ops *MockPrivOps) StartVPN(name string) error {
	ops.startOp()
	defer ops.endOp()
	vpn := ops.mustGetVpn(name)
	if vpn.running {
		panic(fmt.Sprintf("Tried to start already-running vpn %q", name))
	}
	vpn.running = true
	return nil
}

func (ops *MockPrivOps) StopVPN(name string) error {
	ops.startOp()
	defer ops.endOp()
	vpn := ops.mustGetVpn(name)
	if !vpn.running {
		panic(fmt.Sprintf("Tried to stop vpn %q, which is not running.", name))
	}
	vpn.running = false
	return nil
}

func (ops *MockPrivOps) DeleteVPN(name string) error {
	ops.startOp()
	defer ops.endOp()
	vpn, ok := ops.vpns[name]
	if !ok {
		panic(fmt.Sprintf("Tried to delete non-existent vpn %q", name))
	}
	if vpn.running {
		panic(fmt.Sprintf("Tried to delete a running vpn (%q)", name))
	}
	delete(ops.vpns, name)
	return nil
}

func (ops *MockPrivOps) ListVPNs() ([]string, error) {
	ops.startOp()
	defer ops.endOp()
	ret := make([]string, 0, len(ops.vpns))
	for k, _ := range ops.vpns {
		ret = append(ret, k)
	}
	return ret, nil
}

//// Internal consistency stuff.

// Call this at the start of every privileged operation; it locks the ops
// value and does a consistency check.
func (ops *MockPrivOps) startOp() {
	ops.lock.Lock()
	ops.validate()
}

// Call this at the end of every privileged operation; it does a consistency
// check and unlocks the ops value.
func (ops *MockPrivOps) endOp() {
	ops.validate()
	ops.lock.Unlock()
}

// Sanity check the state; panics if any of our invariants are violated.
func (ops *MockPrivOps) validate() {
	// Keep track of what port numbers we've seen as we walk through
	// the set of vpns.
	usedPorts := make(map[uint16]struct{})

	for k, v := range ops.vpns {
		// Make sure the port number isn't used by any vpn we've seen in
		// the past.
		if _, ok := usedPorts[v.portNo]; ok {
			panic(fmt.Sprintf(
				"Port number %d is used by more than one vpn!",
				v.portNo,
			))
		}
		// OK, we're good. Add this to the list for later checks:
		usedPorts[v.portNo] = struct{}{}

		// Make sure the vlan number is valid:
		if v.vlanNo == 0 || v.vlanNo > 4096 {
			panic(fmt.Sprintf(
				"Illegal vlan number for vpn %q: %d",
				k,
				v.vlanNo,
			))
		}
	}
}
