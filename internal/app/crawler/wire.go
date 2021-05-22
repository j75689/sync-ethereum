//+build wireinject

//The build tag makes sure the stub is not built in the final build.

package crawler

import (
	"sync-ethereum/internal/config"
	"sync-ethereum/internal/delivery/crawler"
	crawlerSvc "sync-ethereum/internal/service/ethclient_crawler"
	"sync-ethereum/internal/wireset"

	"github.com/google/wire"
)

func Initialize(configPath string) (Application, error) {
	wire.Build(
		newApplication,
		config.NewConfig,
		wireset.InitLogger,
		wireset.InitMQ,
		crawlerSvc.NewEthClientCrawlerService,
		crawler.NewCrawler,
	)
	return Application{}, nil
}
