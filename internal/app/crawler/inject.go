package crawler

import (
	"sync-ethereum/internal/delivery/crawler"

	"github.com/rs/zerolog"
)

type Application struct {
	logger  zerolog.Logger
	crawler *crawler.Crawler
}

func (application Application) Start() error {
	application.logger.Info().Msg("crawler startup")
	return application.crawler.Start()
}

func (application Application) Stop() error {
	application.logger.Info().Msg("shutdown crawler ...")
	defer application.logger.Info().Msg("crawler is closed")
	return application.crawler.Shutdown()
}

func newApplication(
	logger zerolog.Logger,
	crawler *crawler.Crawler,
) Application {
	return Application{
		logger:  logger,
		crawler: crawler,
	}
}
