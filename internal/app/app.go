package app

import (
	"log/slog"
	"time"

	grpcapp "sso/internal/app/grpc"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(
	log *slog.Logger,
	gRPCPort int,
	storagePath string,
	tokenTTL time.Duration,
) *App {
	// TODO: init storage

	// TODO: init auth service

	grpcApp := grpcapp.New(log, gRPCPort)

	return &App{
		GRPCServer: grpcApp,
	}

}
