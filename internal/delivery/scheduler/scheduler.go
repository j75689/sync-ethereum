package scheduler

import (
	"context"
	"math/big"
	"sync-ethereum/internal/config"
	"sync-ethereum/internal/model"
	"sync-ethereum/internal/service"
	"sync-ethereum/pkg/mq"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

func NewScheduler(config config.Config, logger zerolog.Logger, mq mq.MQ, crawler service.CrawlerService, storageSvc service.StorageService) *Scheduler {
	return &Scheduler{
		config:     config,
		logger:     logger,
		mq:         mq,
		crawler:    crawler,
		storageSvc: storageSvc,
	}
}

type Scheduler struct {
	config     config.Config
	logger     zerolog.Logger
	mq         mq.MQ
	close      func()
	crawler    service.CrawlerService
	storageSvc service.StorageService
}

func (scheduler *Scheduler) Start() error {
	done := make(chan struct{})
	scheduler.close = func() {
		close(done)
	}
	tick := time.NewTicker(scheduler.config.Scheduler.Sync.Interval)
	for {
		select {
		case <-tick.C:
			ctx, cancel := context.WithTimeout(context.Background(), scheduler.config.Scheduler.Sync.Interval)

			number, err := scheduler.crawler.GetBlockNumber(ctx)
			if err != nil {
				scheduler.logger.Error().Err(err).Msg("parse current block number error")
				cancel()
				continue
			}
			scheduler.logger.Info().Msgf("parse current block number: %d", number.Int64())
			onlineBockNumber := model.GormBigInt(*number)

			currentBlockNumber, err := scheduler.storageSvc.GetCurrentBlockNumber(ctx)
			if err != nil {
				scheduler.logger.Error().Err(err).Msg("get database current block number error")
				cancel()
				continue
			}

			if currentBlockNumber.Int64() < scheduler.config.Scheduler.StartAt {
				bi := big.NewInt(scheduler.config.Scheduler.StartAt)
				currentBlockNumber = model.GormBigInt(*bi)

				err = scheduler.storageSvc.UpdateCurrentBlockNumber(ctx, currentBlockNumber, onlineBockNumber)
				if err != nil {
					scheduler.logger.Error().Int64("block_number", currentBlockNumber.Int64()).Err(err).Msg("update db current block number error")
					cancel()
					continue
				}
			}

			i := currentBlockNumber.Int64() - int64(scheduler.config.Scheduler.UnstableNumber) // update unstable block
			limit := i + scheduler.config.Scheduler.BatchLimit
			for i < number.Int64() && i < limit {
				scheduler.logger.Info().Int64("block_number", i).Err(err).Msg("push crawler id")
				if err := scheduler.mq.Publish(scheduler.config.Crawler.Topic, uuid.New().String(), big.NewInt(i).Bytes()); err != nil {
					scheduler.logger.Error().Int64("block_number", i).Err(err).Msg("push crawler id error")
					break
				}
				i++
			}
			number = big.NewInt(i)
			err = scheduler.storageSvc.UpdateCurrentBlockNumber(ctx, model.GormBigInt(*number), onlineBockNumber)
			if err != nil {
				scheduler.logger.Error().Int64("block_number", i).Err(err).Msg("update db current block number error")
				cancel()
				continue
			}

			cancel()
		case <-done:
			return nil
		}
	}
}

func (scheduler *Scheduler) Shutdown() error {
	scheduler.close()
	scheduler.crawler.Close()
	if err := scheduler.storageSvc.Close(); err != nil {
		return err
	}
	return scheduler.mq.Close()
}
