package main

import (
	"context"
	"log"
	"os"
	"tracker-backend/internal/app"
	"tracker-backend/internal/app/dependencies"
	"tracker-backend/internal/app/repository"
	"tracker-backend/internal/config"
	"tracker-backend/internal/pkg/storage"
)

func main() {
	config.MustLoadConfig()

	// TODO: init logger (slog)

	// create empty context
	ctx := context.TODO()

	// create mongodb connection
	mongoClient, err := storage.NewMongoClient(
		ctx,
		os.Getenv(config.MongoURLEnvName),
		os.Getenv(config.MongoDBNameEnvName),
	)
	if err != nil {
		log.Fatalf("failed to connect mongodb %s", err.Error())
	}
	defer mongoClient.Disconnect(ctx)

	// create repository
	repo := repository.MustInitRepository(ctx, mongoClient.Database)

	// create redis connection
	redisClient, err := storage.NewRedisClient(
		ctx,
		os.Getenv(config.RedisHostEnvName),
		os.Getenv(config.RedisPortEnvName),
	)
	if err != nil {
		log.Fatalf("failed to connect redis %s", err.Error())
	}
	defer redisClient.Close()

	// init dependencies
	deps := dependencies.InitDependencies(ctx, repo)

	// create app instance
	app := app.NewApp(
		os.Getenv(config.PortEnvName), deps,
	)

	app.Run()
}
