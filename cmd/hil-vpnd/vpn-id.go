package main

import (
	"encoding/hex"
	"fmt"
	"regexp"
	"strconv"
)

var vpnIdRegexp = regexp.MustCompile("^([0-9]+)-([a-fA-F0-9]+)")

// A VpnId is the name of a vpn as hil-vpnd understands them; we encode
// these as textual names when passing them to the privop command.
type VpnId struct {
	Unique [128 / 8]byte
	Port   uint16
}

// dieBug panics with a message reporting that there is a bug in the software;
// this is used for things which should be impossible.
func dieBug(msg string) {
	panic(fmt.Sprintf("BUG: %s", msg))
}

func (id *VpnId) MarshalText() (text []byte, err error) {
	return []byte(fmt.Sprintf("%d-%x", id.Port, id.Unique)), nil
}

func (id *VpnId) UnmarshalText(text []byte) error {
	ret := vpnIdRegexp.FindSubmatch(text)
	if ret == nil {
		return fmt.Errorf("Could not parse vpn id")
	}
	if len(ret) != 3 {
		dieBug(fmt.Sprintf(
			"Unexpected return size from FindSubmatch: "+
				"%v (should always be 3 on success).",
			ret))
	}
	portNo, err := strconv.ParseUint(string(ret[1]), 10, 16)
	if err != nil {
		dieBug(err.Error())
	}
	id.Port = uint16(portNo)

	unique := make([]byte, len(id.Unique))
	_, err = hex.Decode(unique, ret[2])
	if err != nil {
		dieBug("Error scanning vpn id: " + err.Error())
	}
	copy(id.Unique[:], unique)
	return nil
}
