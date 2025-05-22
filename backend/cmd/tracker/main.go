package main

import (
	"context"
	"log"
	"tracker-backend/internal/app"
	"tracker-backend/internal/app/setup"
	"tracker-backend/internal/config"
	"tracker-backend/internal/pkg/repository"
)

func main() {
	cfg := config.MustLoadConfig()

	// TODO: init logger (slog)

	// create empty context
	ctx := context.TODO()

	// create mongodb connection
	mongoClient, err := repository.NewMongoClient(ctx, cfg.MongoURL, cfg.MongoDBName)
	if err != nil {
		log.Fatalf("failed to connect mongodb %s", err.Error())
	}
	defer mongoClient.Disconnect(ctx)

	// create redis connection
	redisClient, err := repository.NewRedisClient(ctx, cfg.RedisHost, cfg.RedisPort)
	if err != nil {
		log.Fatalf("failed to connect redis %s", err.Error())
	}
	defer redisClient.Close()

	// init dependencies
	deps := setup.InitDependencies(ctx, mongoClient.Database, cfg)

	// create app instance
	app := app.NewApp(cfg.Port, deps)

	app.Run()
}
