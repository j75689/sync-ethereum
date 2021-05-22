package migration

import (
	"sync-ethereum/internal/repository"

	"github.com/rs/zerolog"
)

type Application struct {
	logger zerolog.Logger
	repo   repository.StorageRepository
}

func (application Application) MigrateUpTo(version string) error {
	if err := application.repo.MigrateUpTo(version); err != nil {
		return err
	}
	application.logger.Info().Msgf("migration up to [%s] complete", version)
	return nil
}

func (application Application) MigrateUp() error {
	if err := application.repo.MigrateUp(); err != nil {
		return err
	}
	application.logger.Info().Msg("migration up complete")
	return nil
}

func (application Application) MigrateDownTo(version string) error {
	if err := application.repo.MigrateDownTo(version); err != nil {
		return err
	}
	application.logger.Info().Msgf("migration down to [%s] complete", version)
	return nil
}

func (application Application) MigrateDown() error {
	if err := application.repo.MigrateDown(); err != nil {
		return err
	}
	application.logger.Info().Msg("migration down complete")
	return nil
}

func (application Application) Stop() error {
	return application.repo.Close()
}

func newApplication(
	logger zerolog.Logger,
	repo repository.StorageRepository,
) Application {
	return Application{
		logger: logger,
		repo:   repo,
	}
}
