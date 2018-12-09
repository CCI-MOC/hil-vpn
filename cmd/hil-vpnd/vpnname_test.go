package main

import (
	"testing"
	"testing/quick"
)

// Verify that parseVpnName successfully reverses the output of makeVpnName.
func TestVpnName(t *testing.T) {
	err := quick.Check(func(id UniqueId, port uint16) bool {
		newId, newPort, err := parseVpnName(makeVpnName(id, port))
		ok := err == nil &&
			newId == id &&
			newPort == port
		if !ok {
			t.Logf("Failed; got: %x %x %v", newId, newPort, err)
		}
		return ok
	}, nil)
	if err != nil {
		t.Fatal(err)
	}
}
