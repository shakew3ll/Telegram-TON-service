package appgrpc

import (
	"github.com/shakew3ll/Telegram-TON-service.git/config"
	"github.com/shakew3ll/Telegram-TON-service.git/pkg/logging"
)

type Application struct {
	GRPCSrv *AppGRPC
}

func New(
	logger *logging.Logger,
	cfg config.Config,
) *Application {
	logger.Info("Initializing repositories...")

	logger.Info("Initializing services...")

	grpcApp := NewAppGRPC(logger, cfg.GRPC)

	return &Application{
		GRPCSrv: grpcApp,
	}
}
