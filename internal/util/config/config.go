package config

import (
	"equiprent/internal/util/flags"
	"os"

	yaml "gopkg.in/yaml.v2"
)

// Initialize config
func Initialize() {
	if *flags.Dev {
		Conf = defaultConfig()
	}
	if *flags.Config != "" {
		b, err := os.ReadFile(*flags.Config)
		if err != nil {
			panic(err.Error())
		}
		if err = yaml.Unmarshal(b, &Conf); err != nil {
			panic(err.Error())
		}
	}
}

// Conf is the configuration instance for the running instance
var Conf Config

// defaultConfig initializes the default configuration, when the dev flag is
func defaultConfig() (conf Config) {
	conf.DB.UserENV = "CACHEGRAB_POSTGRESUSER"
	conf.DB.PassENV = "CACHEGRAB_POSTGRESPASS"
	conf.DB.DB = "postgres"
	conf.DB.SSL = !*flags.Dev
	conf.DB.Port = 5432
	conf.DB.Host = "db"
	return
}

// Config is the confugration for the proxyManager application
type Config struct {
	DB DBConfig `yaml:"manager,omitempty"`
}

// DBConfig is the configuration for the manager package
type DBConfig struct {
	UserENV     string `yaml:"userENV,omitempty"`
	LogLocation string `yaml:"logLocation,omitempty"`
	PassENV     string `yaml:"passENV,omitempty"`
	Host        string `yaml:"host,omitempty"`
	Port        int    `yaml:"port,omitempty"`
	DB          string `yaml:"db,omitempty"`
	SSL         bool   `yaml:"ssl,omitempty"`
}
