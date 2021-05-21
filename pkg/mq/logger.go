package mq

import (
	"github.com/ThreeDotsLabs/watermill"
	"github.com/rs/zerolog"
)

var _ watermill.LoggerAdapter = (*WatermillLogger)(nil)

type WatermillLogger struct {
	logger zerolog.Logger
}

func (logger WatermillLogger) Error(msg string, err error, fields watermill.LogFields) {
	logger.logger.Error().Err(err).Fields(fields).Msg(msg)
}

func (logger WatermillLogger) Info(msg string, fields watermill.LogFields) {
	logger.logger.Info().Fields(fields).Msg(msg)
}

func (logger WatermillLogger) Debug(msg string, fields watermill.LogFields) {
	logger.logger.Debug().Fields(fields).Msg(msg)
}

func (logger WatermillLogger) Trace(msg string, fields watermill.LogFields) {
	logger.logger.Trace().Fields(fields).Msg(msg)
}

func (logger WatermillLogger) With(fields watermill.LogFields) watermill.LoggerAdapter {
	return logger
}

func WrapWatermillLogger(logger zerolog.Logger) WatermillLogger {
	return WatermillLogger{logger}
}
