package config

import (
	"backend/models"
	"log"
)

// AppConfig holds the application config i.e. global
type Application struct {
	Config Config
	Logger *log.Logger
	Models models.Models
}

type Config struct {
	Port int
	Env  string
	Db   struct {
		Driver string
		Dsn    string
	}
	Jwt struct {
		Secret string
	}
	Version string
}
