package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

// Тип конфиг соответствует config.yaml
type Config struct {
	Env              string `yaml:"env" env:"ENV" env-default:"local"`
	StoragePath      string `yaml:"storage_path" env-required:"true"`
	ConnectionString string `yaml:"connection_string"`
	HTTPServer       `yaml:"http_server"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

// Функция с приставкой must не должна возрващать значение, а паниковать
func MustLoad() *Config {
	os.Setenv("CONFIG_PATH", `.\cmd\config\local.yaml`)
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg
}
