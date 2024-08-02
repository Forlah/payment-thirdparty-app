package environment

import "os"

type Config struct {
	DatabaseURI  string
	DatabaseName string
	PORT         string
}

func LoadConfig() *Config {
	return &Config{
		DatabaseURI:  os.Getenv("DB_URI"),
		DatabaseName: os.Getenv("DB_NAME"),
		PORT:         os.Getenv("PORT"),
	}
}
