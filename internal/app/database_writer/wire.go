//+build wireinject

//The build tag makes sure the stub is not built in the final build.

package database_writer

import (
	"sync-ethereum/internal/config"
	"sync-ethereum/internal/delivery/database_writer"
	"sync-ethereum/internal/repository/gorm"
	"sync-ethereum/internal/service/storage"
	"sync-ethereum/internal/wireset"

	"github.com/google/wire"
)

func Initialize(configPath string) (Application, error) {
	wire.Build(
		newApplication,
		config.NewConfig,
		wireset.InitLogger,
		wireset.InitDatabase,
		wireset.InitMQ,
		gorm.NewStorageRepository,
		storage.NewStorageService,
		database_writer.NewDatabaseWriter,
	)
	return Application{}, nil
}
