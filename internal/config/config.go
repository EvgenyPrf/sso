package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

type Config struct {
	Env         string        `yaml:"env" env-default:"local"`
	StoragePath string        `yaml:"storage_path" env-required:"true"`
	TokenTTL    time.Duration `yaml:"token_ttl" env-required:"true"`
	GRPC        GRPCConfig    `yaml:"grpc"`
}

type GRPCConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

func MustLoad() *Config {
	path := fetchConfigPath()

	if path == "" {
		panic("config path is empty!")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file does not exist: " + path)
	}

	var config Config

	if err := cleanenv.ReadConfig(path, &config); err != nil {
		panic("failed to read config " + err.Error())
	}

	return &config
}

// получение информации о пути файла конфига из двух возможных источников: либо из переменных окружения, либо из флага
// Priority: flag > env > default.
// Default value is empty string.
func fetchConfigPath() string {
	var res string

	//указание пути файлв конфигурации при запуске приложения
	//пример: --config="path/to/config.yaml"
	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
