package main

import (
	"testing"
)

func TestVpnStates(t *testing.T) {
	states := newStates(config{
		MinPort: 4000,
		MaxPort: 4003,
	})

	vpns := []UniqueId{}

	// Allocate networks, making sure we get ports in the expected order.
	for _, expectedPort := range []uint16{4003, 4002, 4001, 4000} {
		id, actualPort, err := states.NewVpn()
		if err != nil {
			t.Fatal(err)
		}

		if actualPort != expectedPort {
			t.Fatalf("Expected port #%d but got %d.", expectedPort, actualPort)
		}
		vpns = append(vpns, id)
	}

	// Make sure each of the ids was unique.
	for i := range vpns {
		for j := 0; j < i; j++ {
			if vpns[i] == vpns[j] {
				t.Fatalf("vpn Id %x (#%d) is not unique.", vpns[i], i)
			}
		}
	}

	// We should be out of ports now:
	_, _, err := states.NewVpn()
	if err != ErrNoFreePorts {
		t.Fatal("Should have gotten ErrNoFreePorts, but err was ", err)
	}

	// Delete a network, and make sure that we get its port back when
	// we allocate again:

	port, err := states.DeleteVpn(vpns[2])
	if err != nil {
		t.Fatal("Error deleting vpn:", err)
	}
	states.ReleasePort(port)

	_, port, err = states.NewVpn()
	if err != nil {
		t.Fatal("Error allocating vpn ", err)
	}
	if port != 4001 {
		t.Fatalf("Unexpected port number; wanted 4001 but got %d.", port)
	}
}
