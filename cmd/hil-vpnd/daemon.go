package main

import (
	"fmt"
	"net/http"
)

// A Daemon manages the runtime state and configuration of hil-vpnd.
type Daemon struct {
	handler   http.Handler
	privops   PrivOps
	vpnStates *VpnStates
}

// Generate a new daemon using the given config and PrivOps
func newDaemon(cfg config, privops PrivOps) (*Daemon, error) {
	vpnNames, err := privops.ListVPNs()
	if err != nil {
		return nil, fmt.Errorf("Listing existing vpns: %v", err)
	}
	vpnStates := newStates(cfg, vpnNames)

	return &Daemon{
		handler:   makeHandler(cfg.AdminToken, privops, vpnStates),
		privops:   privops,
		vpnStates: vpnStates,
	}, nil
}
