package database_writer

import (
	"context"
	"encoding/json"
	"sync-ethereum/internal/config"
	"sync-ethereum/internal/model"
	"sync-ethereum/internal/service"
	"sync-ethereum/pkg/mq"

	"github.com/rs/zerolog"
)

func NewDatabaseWriter(config config.Config, logger zerolog.Logger, mq mq.MQ, storageSvc service.StorageService) *DatabaseWriter {
	return &DatabaseWriter{
		config:     config,
		logger:     logger,
		mq:         mq,
		storageSvc: storageSvc,
	}
}

type DatabaseWriter struct {
	config     config.Config
	logger     zerolog.Logger
	mq         mq.MQ
	storageSvc service.StorageService
}

func (w *DatabaseWriter) Start() error {
	return w.mq.Subscribe(context.Background(), w.config.DatabaseWriter.PoolSize, w.config.DatabaseWriter.Topic, func(key string, data []byte) (bool, error) {
		ctx, cancel := context.WithTimeout(context.Background(), w.config.DatabaseWriter.Timeout)
		defer cancel()
		block := model.Block{}
		err := json.Unmarshal(data, &block)
		if err != nil {
			return true, err
		}
		err = w.storageSvc.CreateBlock(ctx, &block)
		if err != nil {
			return false, err
		}
		return true, nil
	}, func(key string, e error) {
		w.logger.Error().Str("message_key", key).Err(e).Msg("DatabaseWriter error")
	})
}

func (w *DatabaseWriter) Shutdown() error {
	if err := w.storageSvc.Close(); err != nil {
		return err
	}
	return w.mq.Close()
}
