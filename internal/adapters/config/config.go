package config

import "os"

type (
	App struct {
		Name string
		Env  string
	}
	HTTP struct {
		Env  string
		URL  string
		Port string
	}

	DB struct {
		Uri      string
		Host     string
		Port     string
		User     string
		Password string
		Name     string
	}
	Container struct {
		App  *App
		DB   *DB
		HTTP *HTTP
	}
)

func New() (*Container, error) {
	app := &App{
		Name: os.Getenv("APP_NAME"),
		Env:  os.Getenv("APP_ENV"),
	}

	db := &DB{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Name:     os.Getenv("DB_NAME"),
		Uri:      os.Getenv("DB_URI"),
	}

	http := &HTTP{
		URL:  os.Getenv("APP_URL"),
		Port: os.Getenv("APP_PORT"),
	}

	return &Container{
		app,
		db,
		http,
	}, nil
}
