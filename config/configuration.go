package config

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"time"

	"github.com/heetch/confita"
	"github.com/heetch/confita/backend"
	"github.com/heetch/confita/backend/env"
	"github.com/heetch/confita/backend/flags"
)

// Config is the exporter CLI configuration.
type Config struct {
	AdguardProtocol string        `config:"adguard_protocol"`
	AdguardHostname string        `config:"adguard_hostname"`
	AdguardUsername string        `config:"adguard_username"`
	AdguardPassword string        `config:"adguard_password"`
	AdguardPort     string        `config:"adguard_port"`
	ServerPort      string        `config:"server_port"`
	Interval        time.Duration `config:"interval"`
    LogLimit        string        `config:"log_limit"`
}

func getDefaultConfig() *Config {
	return &Config{
		AdguardProtocol: "http",
		AdguardHostname: "127.0.0.1",
		AdguardUsername: "",
		AdguardPassword: "",
		AdguardPort:     "",
		ServerPort:      "9617",
		Interval:        10 * time.Second,
    LogLimit:        "1000",
	}
}

// Load method loads the configuration by using both flag or environment variables.
func Load() *Config {
	loaders := []backend.Backend{
		env.NewBackend(),
		flags.NewBackend(),
	}

	loader := confita.NewLoader(loaders...)

	cfg := getDefaultConfig()
	err := loader.Load(context.Background(), cfg)
	if err != nil {
		panic(err)
	}

	cfg.show()

	return cfg
}

func (c Config) show() {
	val := reflect.ValueOf(&c).Elem()
	log.Println("---------------------------------------")
	log.Println("- AdGuard Home exporter configuration -")
	log.Println("---------------------------------------")
	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)

		// Do not print password or api token but do print the authentication method
		if typeField.Name != "AdguardPassword" {
			log.Println(fmt.Sprintf("%s : %v", typeField.Name, valueField.Interface()))
		} else {
			showAuthenticationMethod(typeField.Name, valueField.String())
		}
	}
	log.Println("---------------------------------------")
}

func showAuthenticationMethod(name, value string) {
	if len(value) > 0 {
		log.Println(fmt.Sprintf("AdGuard Authentication Method : %s", name))
	}
}
