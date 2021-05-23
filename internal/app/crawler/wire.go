//+build wireinject

//The build tag makes sure the stub is not built in the final build.

package crawler

import (
	"sync-ethereum/internal/config"
	"sync-ethereum/internal/delivery/crawler"
	"sync-ethereum/internal/repository/gorm"
	crawlerSvc "sync-ethereum/internal/service/ethclient_crawler"
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
		crawlerSvc.NewEthClientCrawlerService,
		crawler.NewCrawler,
	)
	return Application{}, nil
}
