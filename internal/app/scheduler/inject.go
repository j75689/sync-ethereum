package scheduler

import (
	"sync-ethereum/internal/delivery/scheduler"

	"github.com/rs/zerolog"
)

type Application struct {
	logger    zerolog.Logger
	scheduler *scheduler.Scheduler
}

func (application Application) Start() error {
	application.logger.Info().Msg("scheduler startup")
	return application.scheduler.Start()
}

func (application Application) Stop() error {
	application.logger.Info().Msg("shutdown scheduler ...")
	defer application.logger.Info().Msg("scheduler is closed")
	return application.scheduler.Shutdown()
}

func newApplication(
	logger zerolog.Logger,
	scheduler *scheduler.Scheduler,
) Application {
	return Application{
		logger:    logger,
		scheduler: scheduler,
	}
}
