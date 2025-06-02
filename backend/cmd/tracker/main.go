package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"time"
	"tracker-backend/internal/config"
	"tracker-backend/internal/domain/repository"
	"tracker-backend/internal/infrastructure/logger"
	"tracker-backend/internal/infrastructure/mongo"
	"tracker-backend/internal/infrastructure/redis"
	"tracker-backend/internal/interfaces/rest"

	"github.com/go-chi/chi/v5"
)

func main() {
	slog.Info("init app...")

	// init .env config
	if err := config.MustLoadConfig(); err != nil {
		panic("could not init config")
	}

	// init logger
	if err := logger.InitLogger(); err != nil {
		panic("could not init logger")
	}

	// connect mongo
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	mongoClient, err := mongo.NewMongoClient(ctx, os.Getenv(config.MongoURLEnvName), os.Getenv(config.MongoDatabaseEnvName))
	if err != nil {
		slog.Error("could not connect mongod", slog.String("error", err.Error()))
	}
	defer mongoClient.Disconnect(ctx)

	// connect redis
	redisClient, err := redis.NewRedisClient(
		ctx,
		os.Getenv(config.RedisHostEnvName),
		os.Getenv(config.RedisPortEnvName),
	)
	if err != nil {
		slog.Error("could not connect redis", slog.String("error", err.Error()))
	}
	_ = redisClient

	// create domain repository
	repo, err := repository.InitRepository(context.TODO(), mongoClient.Database)
	if err != nil {
		slog.Error("could not create repository", slog.String("error", err.Error()))
	}
	_ = repo
	// TODO: init services
	// TODO: init cron jobs

	// create router
	r := chi.NewRouter()
	r.Mount("/api", rest.NewAppRouter())

	// configure server
	server := http.Server{
		Handler: r,
		Addr:    ":" + os.Getenv(config.PortEnvName),
	}

	// starting server
	slog.Info("starting server...", slog.String("port", os.Getenv(config.PortEnvName)))
	if err := server.ListenAndServe(); err != nil {
		slog.Error("error occurred while running server", slog.String("error", err.Error()))
	}
}
