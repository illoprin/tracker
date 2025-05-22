package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func MustLoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("failed to load .env")
	}

	cfg := &Config{}

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatalf("unable to load port")
	}

	mongoUrl := os.Getenv("MONGO_URL")
	if mongoUrl == "" {
		log.Fatalf("unable to load mongodb url")
	}

	mongoDBName := os.Getenv("MONGO_DB_NAME")
	if mongoDBName == "" {
		log.Fatalf("unable to load database name")
	}

	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		log.Fatalf("unable to load redis host")
	}

	redisPort := os.Getenv("REDIS_PORT")
	if redisPort == "" {
		log.Fatalf("unable to load redis port")
	}

	secret := os.Getenv("SECRET")
	if secret == "" {
		log.Fatalf("unable to load jwt secret")
	}

	cfg.MongoURL = mongoUrl
	cfg.RedisHost = redisHost
	cfg.RedisPort = redisPort
	cfg.MongoDBName = mongoDBName
	cfg.Port = port
	cfg.JWTSecret = secret

	return cfg
}
