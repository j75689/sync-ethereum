package crawler

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"sync-ethereum/internal/config"
	"sync-ethereum/internal/model"
	"sync-ethereum/internal/service"
	"sync-ethereum/pkg/mq"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog"
)

func NewCrawler(config config.Config, logger zerolog.Logger, mq mq.MQ, crawler service.CrawlerService) *Crawler {
	return &Crawler{
		config:  config,
		logger:  logger,
		mq:      mq,
		crawler: crawler,
	}
}

type Crawler struct {
	config  config.Config
	logger  zerolog.Logger
	mq      mq.MQ
	crawler service.CrawlerService
}

func (c *Crawler) Start() error {
	return c.mq.Subscribe(context.Background(), c.config.Crawler.PoolSize, c.config.Crawler.Topic, func(key string, data []byte) (bool, error) {
		ctx, cancel := context.WithTimeout(context.Background(), c.config.Crawler.Timeout)
		defer cancel()
		number := new(big.Int).SetBytes(data)
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
				fmt.Println(tx.Hash(), err)
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

		c.logger.Info().RawJSON("block", b).Int64("block_number", number.Int64()).Msg("push to database writer")
		err = c.mq.Publish(c.config.DatabaseWriter.Topic, key, b)
		if err != nil {
			return false, err
		}
		return true, nil
	}, func(key string, e error) {
		c.logger.Error().Str("message_key", key).Err(e).Msg("crawler error")
	})
}

func (c *Crawler) Shutdown() error {
	c.crawler.Close()
	return c.mq.Close()
}
