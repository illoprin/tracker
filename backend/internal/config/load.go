package config

import "github.com/joho/godotenv"

func MustLoadConfig() error {
	return godotenv.Load()
}
