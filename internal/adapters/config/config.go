package config

import (
	"fmt"
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
	return loadEnv()
}

func loadEnv() Config {
	path, err := filepath.Abs(os.Getenv("CONFIG_PATH"))
	if err != nil {
		slog.Error("file path:", err.Error())
	}
	configFileName := fmt.Sprintf("%s/%s.yml", path, os.Getenv("APP_ENV"))
	yamlFile, err := os.ReadFile(configFileName)
	if err != nil {
		slog.Error("Yaml file:", err.Error())
	}
	config := Config{
		App:  &App{},
		HTTP: &HTTP{},
		DB:   &DB{},
	}
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		slog.Error("Unmarshal error:", err.Error())
	}
	return config
}
