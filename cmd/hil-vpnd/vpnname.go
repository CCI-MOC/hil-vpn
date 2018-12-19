package main

import (
	"encoding/hex"
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

var (
	vpnNameRegexp = regexp.MustCompile("^hil_vpn_id_([0-9a-f]{32})_port_([0-9]+)$")

	ErrInvalidVpnName = errors.New("Invalid vpn name")
)

// format the vpn name as we will pass it to PrivOps.
func makeVpnName(id UniqueId, port uint16) string {
	return fmt.Sprintf("hil_vpn_id_%x_port_%d", id, port)
}

func parseVpnName(name string) (id UniqueId, port uint16, err error) {
	matches := vpnNameRegexp.FindStringSubmatch(name)
	if len(matches) != 3 {
		return id, port, ErrInvalidVpnName
	}

	_, err = hex.Decode(id[:], []byte(matches[1]))
	if err != nil {
		return id, port, err
	}

	port64, err := strconv.ParseUint(matches[2], 10, 16)
	if err != nil {
		return id, port, err
	}
	port = uint16(port64)

	return id, port, nil
}
