package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"auth/internal/app"
	"auth/internal/config"
	"auth/internal/lib/logger/slogpretty"
	"auth/internal/storage/postgres"
)

const (
	envLocal = "local"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info("starting service")

	dbURL := dbConnURL(cfg)

	log.Debug("Database connection url", slog.String("conn", dbURL))

	storage, err := postgres.New(dbURL)
	if err != nil {
		panic(err)
	}

	application := app.New(log, cfg.GRPC.DockerPort, storage, cfg.Token.TTL, cfg.Token.Secret)

	go application.GRPCServer.MustRun()

	stop := make(chan os.Signal, 3)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-stop

	application.GRPCServer.Stop()

	if err := storage.Close(); err != nil {
		panic(err)
	}

	log.Info("application stopped gracefully")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = setupPrettyLogger()
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}

func setupPrettyLogger() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}

func dbConnURL(cfg *config.Config) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DB.User, cfg.DB.Password, cfg.DB.Host, cfg.DB.DockerPort, cfg.DB.Name)
}
