package database

import (
	"context"
	"time"

	"github.com/rs/zerolog"
	"gorm.io/gorm/logger"
)

var _ logger.Interface = (*GormLogger)(nil)

type GormLogger struct {
	logger zerolog.Logger
}

func (logger GormLogger) LogMode(logger.LogLevel) logger.Interface {
	return logger
}

func (logger GormLogger) Info(ctx context.Context, format string, args ...interface{}) {
	logger.logger.Info().Msgf(format, args...)
}

func (logger GormLogger) Warn(ctx context.Context, format string, args ...interface{}) {
	logger.logger.Warn().Msgf(format, args...)
}

func (logger GormLogger) Error(ctx context.Context, format string, args ...interface{}) {
	logger.logger.Error().Msgf(format, args...)
}

func (logger GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, row := fc()
	logger.logger.Trace().Dur("elapsed", elapsed).Str("sql", sql).Int64("row", row).Err(err).Msg("trace sql")
}

func WarpGormLogger(logger zerolog.Logger) GormLogger {
	return GormLogger{logger}
}
