package crawler

import (
	"context"
	"encoding/json"
	"sync-ethereum/internal/config"
	"sync-ethereum/internal/model"
	"sync-ethereum/internal/service"
	"sync-ethereum/pkg/mq"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog"
)

func NewCrawler(config config.Config, logger zerolog.Logger, mq mq.MQ, storageSvc service.StorageService, crawler service.CrawlerService) *Crawler {
	return &Crawler{
		config:     config,
		logger:     logger,
		mq:         mq,
		storageSvc: storageSvc,
		crawler:    crawler,
	}
}

type Crawler struct {
	config     config.Config
	logger     zerolog.Logger
	mq         mq.MQ
	storageSvc service.StorageService
	crawler    service.CrawlerService
}

func (c *Crawler) Start() error {
	err := c.mq.Subscribe(context.Background(), c.config.Crawler.PoolSize, c.config.Crawler.Topic, func(key string, data []byte) (bool, error) {
		ctx, cancel := context.WithTimeout(context.Background(), c.config.Crawler.Timeout)
		defer cancel()
		crawlerMessage := model.CrawlerMessage{}
		err := json.Unmarshal(data, &crawlerMessage)
		if err != nil {
			return true, err
		}

		number := crawlerMessage.BlockNumber.BigInt()
		c.logger.Info().Int64("block_number", number.Int64()).Msg("parse block")
		block, err := c.crawler.GetBlockByNumber(ctx, number)
		if err != nil {
			return false, err
		}

		modelBlock := model.Block{
			BlockNumber: model.GormBigInt(*block.Number()),
			BlockHash:   block.Hash().Hex(),
			BlockTime:   block.Time(),
			ParentHash:  block.ParentHash().Hex(),
			IsStable:    crawlerMessage.IsStable,
			Transaction: make([]*model.Transaction, block.Transactions().Len()),
		}

		for idx, tx := range block.Body().Transactions {
			receipt, err := c.crawler.GetTransactionReceipt(ctx, tx.Hash())
			if err != nil {
				return false, err
			}

			from := ""
			signer := types.NewEIP155Signer(tx.ChainId())
			sender, err := signer.Sender(tx)
			if err != nil {
				cancel()
				continue
			}
			from = sender.Hex()

			to := ""
			if tx.To() != nil {
				to = tx.To().Hex()
			}
			modelTx := &model.Transaction{
				BlockNumber: model.GormBigInt(*number),
				TXHash:      tx.Hash().Hex(),
				From:        from,
				To:          to,
				Nonce:       tx.Nonce(),
				Value:       model.GormBigInt(*tx.Value()),
				Data:        model.BlockData(tx.Data()),
				Logs:        make([]*model.TransactionLog, len(receipt.Logs)),
			}

			for logIdx, log := range receipt.Logs {
				modelTx.Logs[logIdx] = &model.TransactionLog{
					TXHash: log.TxHash.Hex(),
					Index:  uint64(log.Index),
					Data:   model.BlockData(log.Data),
				}
			}

			modelBlock.Transaction[idx] = modelTx
		}

		b, err := json.Marshal(modelBlock)
		if err != nil {
			return true, err // format error, not retry
		}

		// pre-written
		err = c.storageSvc.CreateBlock(ctx, &model.Block{
			BlockNumber: modelBlock.BlockNumber,
			BlockHash:   modelBlock.BlockHash,
			BlockTime:   modelBlock.BlockTime,
			ParentHash:  modelBlock.ParentHash,
			IsStable:    false,
		})
		if err != nil {
			c.logger.Error().Err(err).Int64("block_number", number.Int64()).Msg("pre-written block error")
		}

		c.logger.Info().RawJSON("block", b).Int64("block_number", number.Int64()).Msg("push to database writer")
		err = c.mq.Publish(c.config.DatabaseWriter.Topic, key, b)
		if err != nil {
			return false, err
		}
		return true, nil
	}, func(key string, e error) {
		c.logger.Error().Str("message_key", key).Err(e).Msg("crawler error")
	})
	return err
}

func (c *Crawler) Shutdown() error {
	if err := c.mq.Close(); err != nil {
		return err
	}
	if err := c.storageSvc.Close(); err != nil {
		return err
	}
	c.crawler.Close()
	return nil
}
