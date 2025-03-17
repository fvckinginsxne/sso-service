package grpcapp

import (
	"fmt"
	"log/slog"
	"net"

	"google.golang.org/grpc"

	"sso/internal/grpc/authgrpc"
)

type GRPCApp struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

func New(log *slog.Logger, port int, auth authgrpc.Auth) *GRPCApp {
	gRPCServer := grpc.NewServer()

	authgrpc.Register(gRPCServer, auth)

	return &GRPCApp{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
	}
}

// MustRun runs gRPC server and panics if any errors occurs
func (a *GRPCApp) MustRun() {
	if err := a.run(); err != nil {
		panic(err)
	}
}

func (a *GRPCApp) run() error {
	const op = "grpcapp.Run"

	log := a.log.With(
		slog.String("op", op),
		slog.Int("port", a.port),
	)

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("grpc server is running", slog.String("addr", l.Addr().String()))

	if err := a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *GRPCApp) Stop() {
	const op = "grpc.Stop"

	a.log.With(slog.String("op", op)).
		Info("grpc server is stopping")

	a.gRPCServer.GracefulStop()
}
