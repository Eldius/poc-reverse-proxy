package config

import (
	"log"
	"testing"

	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigFile("samples/config.yaml")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Using config file:", viper.ConfigFileUsed())
	}
}

func TestLoadRoutes(t *testing.T) {
	if cfg, err := LoadRoutes(); err != nil {
		t.Error("Failed to load configuration from file")
	} else {
		if len(cfg.Routes) != 4 {
			t.Errorf("Should have 2 routes, but has %d", len(cfg.Routes))
		}
	}
}
