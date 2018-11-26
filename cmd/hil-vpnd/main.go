package main

import (
	"github.com/caarlos0/env"
	"log"
	"net/http"
)

type config struct {
	MinPort int `env:"MIN_VPN_PORT,required"`
	MaxPort int `env:"MAX_VPN_PORT,required"`
}

// Parse and validate the config, then return it.
func getConfig() config {
	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatal("Parsing config from environment: ", err)
	}
	if cfg.MinPort > cfg.MaxPort {
		log.Fatal("Config error: MIN_VPN_PORT is greater than MAX_VPN_PORT")
	}
	if cfg.MinPort < 1024 {
		log.Fatalf("MIN_VPN_PORT specifies a privileged port (%d)", cfg.MinPort)
	}
	if cfg.MaxPort >= (1 << 16) {
		log.Fatalf("MAX_VPN_PORT is out of range (%d)", cfg.MaxPort)
	}
	return cfg
}

func main() {
	cfg := getConfig()
	// TODO: query the state of existent vpns and populate this accordingly:
	vpnStates := newStates()
	for i := cfg.MinPort; i <= cfg.MaxPort; i++ {
		vpnStates.FreePorts = append(vpnStates.FreePorts, uint16(i))
	}

	http.Handle("/", makeHandler(PrivOpsCmd{}, vpnStates))
	panic(http.ListenAndServe(":8080", nil))
}
