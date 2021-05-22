//+build wireinject

//The build tag makes sure the stub is not built in the final build.

package migration

import (
	"sync-ethereum/internal/config"
	"sync-ethereum/internal/repository/gorm"
	"sync-ethereum/internal/wireset"

	"github.com/google/wire"
)

func Initialize(configPath string) (Application, error) {
	wire.Build(
		newApplication,
		config.NewConfig,
		wireset.InitLogger,
		wireset.InitDatabase,
		gorm.NewStorageRepository,
	)
	return Application{}, nil
}
