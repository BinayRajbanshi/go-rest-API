package config

import (
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type HttpServer struct {
	Address string `yaml:"address" env:"ENDPOINT"`
}

type Config struct {
	Env          string `yaml:"env" env:"ENV" env-required:"true" `
	DatabasePath string `yaml:"database_path" env:"DATABASE_PATH" `
	HttpServer   `yaml:"http_server"`
}

func MustLoad() *Config {
	var configPath string

	configPath = os.Getenv("CONFIG_PATH")

	// First tries to get the config path from the CONFIG_PATH environment variable
	// If not found, looks for a --config command line flag
	// If neither exists, terminates with a fatal error message

	if configPath == "" {
		flags := flag.String("config", "", "path to the config file")
		flag.Parse()

		configPath = *flags

		if configPath == "" {
			log.Fatal("Config path is not set")
		}
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("Config file does not exist: %s\n", configPath)
	}

	// Verifies the config file exists, reads it, parses it and finally returns the Config struct pointer
	var cfg Config

	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatalf("Cannot read config file: %s", err.Error())
	}

	return &cfg
}

// yaml tags: Define how fields are named in YAML files
// env tags: Define which environment variables can override config values
// env-required: Marks fields that must be set
