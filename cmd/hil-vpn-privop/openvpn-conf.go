package main

import (
	"text/template"
)

// Template for the open vpn config files we generate.
var openVpnCfgTpl = template.Must(template.New("openvpn-config").Parse(`
dev tap-hil-vpn-{{ .Name }}
secret hil-vpn-{{ .Name }}.key
user nobody
group nobody
`))

type OpenVpnCfg struct {
	Name string
	Key  string
}
