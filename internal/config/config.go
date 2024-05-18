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
	Metrics          `yaml:"metrics"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8000"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
	User        string        `yaml:"user" env-required:"true"`
	Password    string        `yaml:"password" env-required:"true" env:"HTTP_SERVER_PASSWORD"`
}

type Metrics struct {
	PrometheusAddress string `yaml:"address" env-default:"localhost:9001"`
}

// MustLoad - функция, которая загружает конфиг
// из указанного в переменной окружения CONFIG_PATH файла
// если CONFIG_PATH не задана, то паниковать
// если файла не существует, то паниковать
// если при чтении конфига возникла ошибка, то паниковать
func MustLoad() *Config {
	// Установим путь до конфига
	os.Setenv("CONFIG_PATH", `.\cmd\config\local.yaml`)

	// Получим путь к конфигу из переменной окружения
	configPath := os.Getenv("CONFIG_PATH")

	// Если путь не задан, то паниковать
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	// Проверим, что файл с конфигом существует
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config

	// Читаем конфиг из файла
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	// Возвращаем ссылку на считанный конфиг
	return &cfg
}
