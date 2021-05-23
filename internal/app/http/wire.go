//+build wireinject

//The build tag makes sure the stub is not built in the final build.

package http

import (
	"github.com/google/wire"
	"sync-ethereum/internal/config"
	"sync-ethereum/internal/delivery/http"
	"sync-ethereum/internal/repository/gorm"
	"sync-ethereum/internal/service/storage"
	"sync-ethereum/internal/wireset"
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
		http.NewHttpServer,
	)
	return Application{}, nil
}
