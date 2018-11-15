package main

// The actual functionality of the 'create' subcommand. The caller is
// responsible for validating the arguments. If an error occurs,
// the program will exit with an error. Otherwise, the generated
// static key is returned.
func createCmd(vpnName string, vlanNo, portNo uint16) string {
	cfg, err := NewOpenVpnConfig(vpnName, portNo)
	chkfatal("Generating openvpn config:", err)
	chkfatal("Saving openvpn config:", cfg.Save())
	return cfg.Key
}
