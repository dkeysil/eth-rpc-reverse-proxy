package config

import (
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Server struct {
		Host string `json:"host" env-description:"Server host" env-default:"localhost"`
		Port string `json:"port" env-description:"Server port" env-default:"8080"`
	} `json:"server"`
	HTTPUpstreams map[string][]string `json:"httpUpstreams" env-description:"List of backends"`
	WSUpstreams   map[string][]string `json:"wsUpstreams" env-description:"List of backends"`
}

func NewConfig() (cfg Config, _ error) {
	return cfg, cleanenv.ReadConfig(os.Getenv("CONFIG_PATH"), &cfg)
}
