package config

import (
	"gopkg.in/yaml.v3"
	"log/slog"
	"os"
	"path/filepath"
)

type (
	App struct {
		Name string `yaml:"name"`
		Env  string `yaml:"environment"`
	}
	HTTP struct {
		GinMode string `yaml:"ginMode"`
		URL     string `yaml:"url"`
		Port    string `yaml:"port"`
	}

	DB struct {
		Uri      string `yaml:"uri"`
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		User     string `yaml:"username"`
		Password string `yaml:"password"`
		Name     string `yaml:"dbName"`
	}

	Config struct {
		App  *App  `yaml:"app"`
		HTTP *HTTP `yaml:"server"`
		DB   *DB   `yaml:"mongo"`
	}
)

func New() Config {
	filePath, err := filepath.Abs(os.Getenv("CONFIG_PATH"))
	if err != nil {
		slog.Error("File path:", err.Error())
	}
	yamlFile, err := os.ReadFile(filePath)
	if err != nil {
		slog.Error("Yaml file:", err.Error())
	}
	var config Config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		slog.Error("Unmarshal error:", err.Error())
	}
	return config
}
