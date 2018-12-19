package main

import (
	"log"
	"net/http"

	"github.com/caarlos0/env"

	"github.com/CCI-MOC/obmd/httpserver"
)

type config struct {
	MinPort      int `env:"MIN_VPN_PORT,required"`
	MaxPort      int `env:"MAX_VPN_PORT,required"`
	ServerConfig httpserver.Config
}

// Parse and validate the config, then return it.
func getConfig() config {
	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatal("Parsing config from environment: ", err)
	}
	if err := env.Parse(&cfg.ServerConfig); err != nil {
		log.Fatal("Parsing config from environment: ", err)
	}
	if err := cfg.ServerConfig.Validate(); err != nil {
		log.Fatal(err)
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
	daemon, err := newDaemon(cfg, PrivOpsCmd{})
	if err != nil {
		log.Fatal(err)
	}
	http.Handle("/", daemon.handler)
	panic(httpserver.Run(&cfg.ServerConfig, nil))
}
