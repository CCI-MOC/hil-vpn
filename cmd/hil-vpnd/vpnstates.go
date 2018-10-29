package main

import (
	"sync"
)

// Track the currently existent vpns, available port numbers, etc.
type VpnStates struct {
	sync.Mutex

	// The ports used by each vpn.
	UsedPorts map[UniqueId]uint16

	// A list of free ports, which may be used with new vpns.
	FreePorts []uint16
}

// Allocate an empty VpnStates.
func newStates() *VpnStates {
	return &VpnStates{
		UsedPorts: map[UniqueId]uint16{},
		FreePorts: []uint16{},
	}
}
