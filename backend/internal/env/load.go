package env

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	MongoURL    string
	MongoDBName string
	RedisURL    string
	JWTSecret   string
}

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

	redisUrl := os.Getenv("REDIS_URL")
	if redisUrl == "" {
		log.Fatalf("unable to load redis url")
	}

	secret := os.Getenv("SECRET")
	if secret == "" {
		log.Fatalf("unable to load jwt secret")
	}

	cfg.MongoURL = mongoUrl
	cfg.RedisURL = redisUrl
	cfg.MongoDBName = mongoDBName
	cfg.Port = port
	cfg.JWTSecret = secret

	return cfg
}
