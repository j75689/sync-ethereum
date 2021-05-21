package wireset

import (
	"sync-ethereum/internal/config"
	"sync-ethereum/pkg/logger"

	"github.com/rs/zerolog"
)

func InitLogger(config config.Config) (zerolog.Logger, error) {
	return logger.NewLogger(config.Logger.Level, config.Logger.Format, logger.WithStr("app_id", config.APPID))
}
