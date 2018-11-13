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
