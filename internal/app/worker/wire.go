//+build wireinject

//The build tag makes sure the stub is not built in the final build.

package worker

import (
	"sync-ethereum/internal/config"
	"sync-ethereum/internal/delivery/worker"
	"sync-ethereum/internal/repository/gorm"
	crawler "sync-ethereum/internal/service/ethclient_crawler"
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
		crawler.NewEthClientCrawlerService,
		storage.NewStorageService,
		worker.NewWorker,
	)
	return Application{}, nil
}
