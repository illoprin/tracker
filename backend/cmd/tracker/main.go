package main

import (
	"context"
	"fmt"
	"log"
	"tracker-backend/internal/app"
	"tracker-backend/internal/env"
	mng "tracker-backend/internal/lib/repository/mongod"
	"tracker-backend/internal/service"
)

func main() {
	cfg := env.MustLoadConfig()

	// TODO: init logger (slog)

	// create empty context
	ctx := context.TODO()

	mongoClient, err := mng.NewMongoClient(ctx, cfg.MongoURL, cfg.MongoDBName)
	if err != nil {
		log.Fatalf("failed to connect mongodb %s", err.Error())
	}
	defer func() {
		err := mongoClient.Disconnect(ctx)
		if err != nil {
			fmt.Println(err)
		}
	}()

	// TODO: init redis connection

	// init dependencies
	deps := service.InitDependencies(ctx, mongoClient.Database, cfg)

	// create app instance
	app := app.NewApp(cfg.Port, deps)

	app.Run()
}
