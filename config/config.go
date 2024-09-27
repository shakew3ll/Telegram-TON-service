package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	IsDebug *bool         `yaml:"is_debug"`
	GRPC    GRPCConfig    `yaml:"grpc"`
	Logger  LoggerConfig  `yaml:"logger"`
	Timeout TimeoutConfig `yaml:"timeout"`
}

type GRPCConfig struct {
	Port    int           `yaml:"port"`
	Host    string        `yaml:"host"`
	Timeout time.Duration `yaml:"timeout"`
}

type LoggerConfig struct {
	Writer string `yaml:"writer"`
	Level  string `yaml:"level"`
}

type TimeoutConfig struct {
	Value time.Duration `yaml:"value"`
}

func New() (Config, error) {
	instance := Config{}
	log.Println("Reading application's configuration")

	useEnvConfig := os.Getenv("USE_ENV_CONFIG") == "true"

	if !useEnvConfig {
		if err := cleanenv.ReadConfig("config.yml", &instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			log.Println("Configuration file error:")
			log.Println(help)
			return Config{}, err
		}
	}

	if err := cleanenv.ReadEnv(&instance); err != nil {
		return Config{}, err
	}

	return instance, nil
}
