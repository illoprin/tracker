package config

import (
	"github.com/joho/godotenv"
)

// load dotenv
// if could not load - panic
func MustLoadConfig() {
	err := godotenv.Load()
	if err != nil {
		panic("failed to load .env")
	}
}
