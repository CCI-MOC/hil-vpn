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

// NOTE: VpnStates has some methods which are thread safe, and others
// which are not. The exported (capitalized-names) methods lock the
// VpnStates, and so are thread-safe, while the lower-case named
// methods are not (but may be called when the VpnStates is already
// locked).

// Allocate a fresh VpnStates, with `FreePorts` generated based on
// the config. The `vpnNames` argument should be the output of
// PrivOps.ListVPNs.
func newStates(cfg config, vpnNames []string) *VpnStates {
	ret := &VpnStates{
		UsedPorts: map[UniqueId]uint16{},
		FreePorts: []uint16{},
	}
	usedPorts := make(map[uint16]struct{})
	for _, name := range vpnNames {
		id, port, err := parseVpnName(name)
		if err != nil {
			// skip it; perhaps the local sysadmin created an
			// openvpn config unrelated to hil-vpn.
			continue
		}
		ret.UsedPorts[id] = port
		usedPorts[port] = struct{}{}
	}
	for i := cfg.MinPort; i <= cfg.MaxPort; i++ {
		if _, ok := usedPorts[uint16(i)]; ok {
			// Port is in use; leave it out.
			continue
		}
		ret.FreePorts = append(ret.FreePorts, uint16(i))
	}
	return ret
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
//
// Note that this does *not* return the vpn's port to the free
// pool; that must be done separately, via ReleasePort()
func (s *VpnStates) DeleteVpn(id UniqueId) (uint16, error) {
	s.Lock()
	defer s.Unlock()

	portNo, ok := s.UsedPorts[id]
	if !ok {
		return 0, ErrNoSuchVpn
	}
	delete(s.UsedPorts, id)
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
func (s *VpnStates) ReleasePort(portNo uint16) {
	s.Lock()
	defer s.Unlock()

	s.FreePorts = append(s.FreePorts, portNo)
}
