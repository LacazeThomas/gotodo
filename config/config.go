package config

import "os"

type DB struct {
	Dialect  string `env:"Dialect" envDefault:"postgres"`
	Host     string `env:"Host"`
	Port     int    `env:"Port" envDefault:"5432"`
	Username string `env:"Username"`
	Password string `env:"Password"`
	Name     string `env:"Name"`
	Charset  string `env:"Charset" envDefault:"utf8"`
}

func GetTokenString() string {
	return os.Getenv("TokenString")
}
