package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env        string     `env-default:"local" yaml:"env"`
	HTTPServer HTTPServer `yaml:"http_server"`
}

type HTTPServer struct {
	Address     string        `env-default:"localhost:8080" yaml:"address"`
	TimeOut     time.Duration `env-default:"4s" yaml:"timeout"`
	IdleTimeout time.Duration `env-default:"60s" yaml:"idle_timeout"`
}

func MustLoad() *Config {

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	if _, err := os.Stat(configPath); os.IsExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg
}
