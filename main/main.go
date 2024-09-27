package main

import (
	"context"
	"log"
	_ "net/http/pprof"
	"os/signal"
	"syscall"

	"github.com/shakew3ll/Telegram-TON-service.git/config"
	"github.com/shakew3ll/Telegram-TON-service.git/infrastructure/appgrpc"
	"github.com/shakew3ll/Telegram-TON-service.git/pkg/logging"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("Failed to load the config due an error: %v", err.Error())
	}

	logger, err := logging.New(cfg)
	if err != nil {
		logger.Fatalf("Failed to configure logger due an error: %v", err.Error())
	}
	logger.Info("Logger connected successfully.")

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	app := appgrpc.New(logger, cfg)
	go app.GRPCSrv.MustRun()

	<-ctx.Done()
	logger.Info("Received shutdown signal, shutting down application...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.Timeout.Value)
	defer cancel()

	app.GRPCSrv.Stop(shutdownCtx)
}
