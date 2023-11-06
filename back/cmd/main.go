package main

import (
	"context"
	"flag"
	"github.com/aakosarev/kanban-board/back/internal/config"
	"github.com/aakosarev/kanban-board/back/internal/server"
	"github.com/aakosarev/kanban-board/back/pkg/logger"
	"github.com/aakosarev/kanban-board/back/pkg/postgres"
	"github.com/aakosarev/kanban-board/back/pkg/redis"
	"log"
	"time"
)

func main() {
	flag.Parse()

	cfg, err := config.InitConfig()
	if err != nil {
		log.Fatal(err)
	}

	appLogger := logger.NewAppLogger(cfg.Logger)
	appLogger.InitLogger()
	appLogger.WithName(server.GetMicroserviceName(cfg))

	postgresConfig := postgres.NewPgConfig(
		cfg.Postgres.Username,
		cfg.Postgres.Password,
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.Database,
	)

	postgresClient, err := postgres.NewClient(context.Background(), 5, time.Second*5, postgresConfig)
	if err != nil {
		appLogger.Fatalf("Postgres init: %s", err)
	} else {
		appLogger.Info("Postgres connected")
	}

	redisClient := redis.NewRedisClient(cfg)
	defer redisClient.Close()
	appLogger.Info("Redis connected")

	appLogger.Fatal(server.NewServer(cfg, appLogger, redisClient, postgresClient).Run())
}
