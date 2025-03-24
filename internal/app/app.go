package app

import (
	"log/slog"
	"time"

	grpcapp "sso/internal/app/grpc"
	"sso/internal/services/auth"
	"sso/internal/storage/sqlite"
)

type App struct {
	GRPCServer *grpcapp.GRPCApp
}

func New(log *slog.Logger,
	gRPCPort int,
	storagePath string,
	tokenTTL time.Duration,
) *App {
	storage, err := sqlite.New(storagePath)
	if err != nil {
		panic(err)
	}

	authService := auth.New(log, storage, storage, storage, tokenTTL)

	grpcApp := grpcapp.New(log, gRPCPort, authService)

	return &App{
		GRPCServer: grpcApp,
	}
}
