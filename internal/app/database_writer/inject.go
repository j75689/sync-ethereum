package database_writer

import (
	"sync-ethereum/internal/delivery/database_writer"

	"github.com/rs/zerolog"
)

type Application struct {
	logger          zerolog.Logger
	database_writer *database_writer.DatabaseWriter
}

func (application Application) Start() error {
	application.logger.Info().Msg("database_writer startup")
	return application.database_writer.Start()
}

func (application Application) Stop() error {
	application.logger.Info().Msg("shutdown database_writer ...")
	defer application.logger.Info().Msg("database_writer is closed")
	return application.database_writer.Shutdown()
}

func newApplication(
	logger zerolog.Logger,
	database_writer *database_writer.DatabaseWriter,
) Application {
	return Application{
		logger:          logger,
		database_writer: database_writer,
	}
}
