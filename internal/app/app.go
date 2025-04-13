package app

import (
	"log/slog"
	"time"

	"sso/internal/app/grpcapp"
	"sso/internal/services/auth"
	"sso/internal/storage/postgres"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(
	log *slog.Logger,
	grpcPort int,
	dbConnURL string,
	tokenTTL time.Duration,
) *App {
	storage, err := postgres.New(dbConnURL)
	if err != nil {
		panic(err)
	}

	authService := auth.New(log, storage, storage, storage, tokenTTL)

	grpcApp := grpcapp.New(log, authService, grpcPort)

	return &App{
		GRPCServer: grpcApp,
	}
}
