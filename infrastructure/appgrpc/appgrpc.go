package appgrpc

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"

	"github.com/shakew3ll/Telegram-TON-service.git/config"
	"github.com/shakew3ll/Telegram-TON-service.git/pkg/logging"
)

type AppGRPC struct {
	logger     *logging.Logger
	gRPCServer *grpc.Server
	host       string
	port       int
}

func NewAppGRPC(
	logger *logging.Logger,
	cfg config.GRPCConfig,
) *AppGRPC {
	gRPCServer := grpc.NewServer()

	return &AppGRPC{
		logger:     logger,
		gRPCServer: gRPCServer,
		host:       cfg.Host,
		port:       cfg.Port,
	}
}

func (a *AppGRPC) MustRun() {
	defer func() {
		if r := recover(); r != nil {
			a.logger.Fatalf("Application panicked: %v", r)
		}
	}()

	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *AppGRPC) Run() error {
	a.logger.Info("Starting gRPC server...")

	lsn, err := net.Listen("tcp", fmt.Sprintf("%s:%d", a.host, a.port))
	if err != nil {
		a.logger.Fatalf("Failed to start TCP listening due an error: %v", err)
		return err
	}

	a.logger.Infof("gRPC server is listening on host: %s", lsn.Addr().String())

	if err := a.gRPCServer.Serve(lsn); err != nil {
		a.logger.Fatalf("Failed to start gRPC server due an error: %v", err)
		return err
	}

	return nil
}

func (a *AppGRPC) Stop(ctx context.Context) {
	done := make(chan struct{})

	go func() {
		defer close(done)
		a.logger.Info("Stopping gRPC server...")
		a.gRPCServer.GracefulStop()
	}()

	select {
	case <-done:
		a.logger.Info("Server gracefully stopped.")
	case <-ctx.Done():
		a.logger.Warn("Graceful shutdown timed out, forcing server stop.")
		a.gRPCServer.Stop()
	}
}
