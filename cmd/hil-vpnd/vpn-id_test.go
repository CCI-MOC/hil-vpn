package main

import (
	"testing"
	"testing/quick"

	"encoding"
)

var (
	// Make sure *VpnId implements these interfaces:
	_ encoding.TextMarshaler   = &VpnId{}
	_ encoding.TextUnmarshaler = &VpnId{}
)

func TestVpnIdMarshalUnmarshal(t *testing.T) {
	err := quick.Check(func(id VpnId) bool {
		text, err := (&id).MarshalText()
		if err != nil {
			t.Fatal(err)
		}

		newId := &VpnId{}
		err = newId.UnmarshalText(text)
		if err != nil {
			t.Fatal(err)
		}
		ret := id == *newId
		if !ret {
			t.Logf("Source was %v, marshalled as %q, unmarshalled as %v",
				id, text, *newId)
		}
		return ret
	}, nil)
	if err != nil {
		t.Fatal(err)
	}
}
