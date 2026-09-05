package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

// @title Junior Project
// @version 1.0

func main() {
	ctx := context.Background()

	// Читаем конфигурацию из переменных окружения
	dbHost := getEnv("DB_HOST", "127.0.0.1")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "Comix")
	dbPassword := getEnv("DB_PASSWORD", "tron")
	dbName := getEnv("DB_NAME", "product-store")
	serverPort := getEnv("SERVER_PORT", "8081")

	// Формируем строку подключения к БД
	dbAddress := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	cfg := config{
		address: ":" + serverPort,
		db: dbConfig{
			address: dbAddress,
		},
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	poolConfig, err := pgxpool.ParseConfig(cfg.db.address)
	if err != nil {
		panic(err)
	}

	poolConfig.MaxConns = 10
	poolConfig.MinConns = 2

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		panic(err)
	}
	defer pool.Close()

	logger.Info("Connected to database", cfg.db.address)

	api := &application{
		config: cfg,
		db:     pool,
	}

	if err := api.run(api.mount()); err != nil {
		slog.Error("server failed to start", "error", err)
		os.Exit(1)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
