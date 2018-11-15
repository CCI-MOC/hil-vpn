package main

import (
	"fmt"
	"os"
	"os/exec"
	"text/template"
)

const configDir = "/etc/openvpn/server"

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

// Save the openvpn config and its static keys to disk.
func (cfg OpenVpnCfg) Save() error {
	cfgPath := configDir + "/" + cfg.Name + ".conf"
	keyPath := configDir + "/hil-vpn-" + cfg.Name + ".key"

	cfgFile, err := os.OpenFile(cfgPath, os.O_CREATE|os.O_EXCL|os.O_RDWR, 0600)
	if err != nil {
		return err
	}
	defer func() {
		cfgFile.Close()
		if err != nil {
			os.Remove(cfgPath)
		}
	}()
	keyFile, err := os.OpenFile(keyPath, os.O_CREATE|os.O_EXCL|os.O_RDWR, 0600)
	if err != nil {
		return err
	}
	defer func() {
		keyFile.Close()
		if err != nil {
			os.Remove(keyPath)
		}
	}()
	if err = openVpnCfgTpl.Execute(cfgFile, cfg); err != nil {
		return err
	}
	_, err = keyFile.Write([]byte(cfg.Key))
	return err
}

// Generate a new openvpn config (including a static key).
func NewOpenVpnConfig(name string) (*OpenVpnCfg, error) {
	cmd := exec.Command("openvpn", "--genkey", "--secret", "/dev/fd/1")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("Error invoking openvpn: %v", err)
	}
	return &OpenVpnCfg{
		Name: name,
		Key:  string(output),
	}, nil
}
