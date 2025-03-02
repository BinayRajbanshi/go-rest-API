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

	var cfg Config

	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatalf("Cannot read config file: %s", err.Error())
	}

	return &cfg
}

// yaml -> take variable from yaml config file
// env -> take variable from .env file and overwrite the value from yaml file.
// env-required -> validates the data
// env-default -> if unable to get value from the files then use this value
