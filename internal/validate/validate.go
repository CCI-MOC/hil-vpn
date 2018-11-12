package validate

import (
	"fmt"
	"regexp"
)

// A regular expression matching legal vpn names.
var vpnNameRegexp = regexp.MustCompile("^[-_a-zA-Z0-9]+$")

// Check whether `name` is a legal name for a vpn. If so, return nil,
// otherwise return an error
func CheckVpnName(name string) error {
	if vpnNameRegexp.MatchString(name) {
		return nil
	}
	return fmt.Errorf("Invalid vpn name %q; names may only contain dashes, underscores, "+
		"and alphanumeric characters.\n\n",
		name)
}

// Check whether `vlanNo` is a valid vlan id. If so, return nil,
// otherwise return an error.
func CheckVlanNo(vlanNo uint16) error {
	if 0 < vlanNo && vlanNo < 4095 {
		return nil
	}
	return fmt.Errorf(
		"Invalid Vlan ID #%d; Vlan IDs must be in the range [1,4094] (inclusive)",
		vlanNo)
}
