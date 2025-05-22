package config

type Config struct {
	Port        string
	MongoURL    string
	MongoDBName string
	RedisHost   string
	RedisPort   string
	JWTSecret   string
}
