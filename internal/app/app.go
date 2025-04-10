package app

import (
	"log/slog"
	"time"

	"sso/internal/app/grpcapp"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(log *slog.Logger,
	grpcPort int,
	dbConnURL string,
	tokenTTL time.Duration,
) *App {
	// TODO: init storage

	// TODO: init auth service

	grpcApp := grpcapp.New(log, grpcPort)

	return &App{
		GRPCServer: grpcApp,
	}
}
