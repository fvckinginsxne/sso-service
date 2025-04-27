package app

import (
	"log/slog"
	"time"

	"auth/internal/app/grpcapp"
	"auth/internal/services/auth"
	"auth/internal/storage/postgres"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(
	log *slog.Logger,
	grpcPort int,
	storage *postgres.Storage,
	tokenTTL time.Duration,
	tokenSecret string,
) *App {
	authService := auth.New(log, storage, storage, tokenTTL, tokenSecret)

	grpcApp := grpcapp.New(log, authService, grpcPort)

	return &App{
		GRPCServer: grpcApp,
	}
}
