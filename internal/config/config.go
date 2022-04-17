package config

import (
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type config struct {
	Server struct {
		Host string `json:"host" env-description:"Server host" env-default:"localhost"`
		Port string `json:"port" env-description:"Server port" env-default:"8080"`
	} `json:"server"`
	Upstreams map[string][]string `json:"upstreams" env-description:"List of backends"`
}

func NewConfig() (cfg config, _ error) {
	return cfg, cleanenv.ReadConfig(os.Getenv("CONFIG_PATH"), &cfg)
}
