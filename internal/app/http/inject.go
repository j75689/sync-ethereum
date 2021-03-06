package http

import (
	"fmt"
	"sync-ethereum/internal/config"
	"sync-ethereum/internal/delivery/http"

	"github.com/rs/zerolog"
)

type Application struct {
	logger     zerolog.Logger
	config     config.Config
	httpServer *http.HttpServer
}

func (application Application) Start() error {
	application.logger.Info().Msgf("http server listen :%d", application.config.HTTP.Port)
	return application.httpServer.Run(fmt.Sprintf(":%d", application.config.HTTP.Port))
}

func (application Application) Stop() error {
	application.logger.Info().Msg("shutdown http server ...")
	defer application.logger.Info().Msg("http server is closed")
	return application.httpServer.Shutdown()
}

func newApplication(
	logger zerolog.Logger,
	config config.Config,
	httpServer *http.HttpServer,
) Application {
	return Application{
		logger:     logger,
		config:     config,
		httpServer: httpServer,
	}
}
