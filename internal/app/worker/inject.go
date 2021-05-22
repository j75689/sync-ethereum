package worker

import (
	"sync-ethereum/internal/delivery/worker"

	"github.com/rs/zerolog"
)

type Application struct {
	logger zerolog.Logger
	worker *worker.Worker
}

func (application Application) Start() error {
	application.logger.Info().Msg("worker startup")
	return application.worker.Start()
}

func (application Application) Stop() error {
	application.logger.Info().Msg("shutdown worker ...")
	defer application.logger.Info().Msg("worker is closed")
	return application.worker.Shutdown()
}

func newApplication(
	logger zerolog.Logger,
	worker *worker.Worker,
) Application {
	return Application{
		logger: logger,
		worker: worker,
	}
}
