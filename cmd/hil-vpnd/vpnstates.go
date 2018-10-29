package main

import (
	"crypto/rand"
	"errors"
	"sync"
)

var (
	// Error indicating that we're out of free ports for openvpn to listen on.
	ErrNoFreePorts = errors.New("There are no free OpenVPN ports")

	// Error indicating that a specified vpn does not exist.
	ErrNoSuchVpn = errors.New("There is no such vpn")
)

// A unique identifier for a vpn.
type UniqueId [128 / 8]byte

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

// Allocate a new vpn. Returns a unique id and a port number.
// May return ErrNoFreePorts if we're out of port numbers to assign.
func (s *VpnStates) NewVpn() (UniqueId, uint16, error) {
	s.Lock()
	defer s.Unlock()

	var id UniqueId
	if _, err := rand.Read(id[:]); err != nil {
		return id, 0, err
	}

	portNo, err := s.allocPort()
	if err == nil {
		s.UsedPorts[id] = portNo
	}
	return id, portNo, err
}

// Delete a vpn. This returns the port number and an error which
// will either be nil or ErrNoSuchVpn.
func (s *VpnStates) DeleteVpn(id UniqueId) (uint16, error) {
	s.Lock()
	defer s.Unlock()

	portNo, ok := s.UsedPorts[id]
	if !ok {
		return 0, ErrNoSuchVpn
	}
	delete(s.UsedPorts, id)
	s.releasePort(portNo)
	return portNo, nil
}

// Allocate a new port for a vpn.
func (s *VpnStates) allocPort() (uint16, error) {
	if len(s.FreePorts) == 0 {
		return 0, ErrNoFreePorts
	}

	portNo := s.FreePorts[len(s.FreePorts)-1]
	s.FreePorts = s.FreePorts[:len(s.FreePorts)-1]
	return portNo, nil
}

// Return a port number to the free pool.
func (s *VpnStates) releasePort(portNo uint16) {
	s.FreePorts = append(s.FreePorts, portNo)
}
