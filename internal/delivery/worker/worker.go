package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"sync-ethereum/internal/config"
	"sync-ethereum/internal/model"
	"sync-ethereum/internal/service"
	"sync-ethereum/pkg/mq"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"gorm.io/gorm"
)

func NewWorker(config config.Config, db *gorm.DB, mq mq.MQ, crawler service.CrawlerService, storageSvc service.StorageService) *Worker {
	return &Worker{
		config:     config,
		db:         db,
		mq:         mq,
		crawler:    crawler,
		storageSvc: storageSvc,
	}
}

type Worker struct {
	config     config.Config
	db         *gorm.DB
	mq         mq.MQ
	close      func()
	crawler    service.CrawlerService
	storageSvc service.StorageService
}

func (worker *Worker) Start() error {
	done := make(chan struct{})
	worker.close = func() {
		close(done)
	}
	tick := time.NewTicker(worker.config.Worker.Sync.Interval)
	for {
		select {
		case <-tick.C:
			ctx, cancel := context.WithTimeout(context.Background(), worker.config.Worker.Sync.Interval)
			number, err := worker.crawler.GetBlockNumber(ctx)
			if err != nil {
				fmt.Println(err)
				cancel()
				continue
			}

			block, err := worker.crawler.GetBlockByNumber(ctx, number)
			if err != nil {
				fmt.Println(err)
				cancel()
				continue
			}
			modelBlock := model.Block{
				BlockNumber: model.GormBigInt(*block.Number()),
				BlockHash:   block.Hash().Hex(),
				BlockTime:   block.Time(),
				ParentHash:  block.ParentHash().Hex(),
				Transaction: make([]*model.Transaction, block.Transactions().Len()),
			}

			for idx, tx := range block.Body().Transactions {
				receipt, err := worker.crawler.GetTransactionReceipt(ctx, tx.Hash())
				if err != nil {
					fmt.Println(tx.Hash(), err)
					cancel()
					continue
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
				fmt.Println(err)
			} else {
				fmt.Println(string(b))
			}

			fmt.Println(worker.storageSvc.CreateBlock(ctx, &modelBlock))
			cancel()
		case <-done:
			return nil
		}
	}
}

func (worker *Worker) Shutdown() error {
	worker.close()
	worker.crawler.Close()
	err := worker.mq.Close()
	if err != nil {
		return err
	}
	db, err := worker.db.DB()
	if err != nil {
		return err
	}
	err = db.Close()
	if err != nil {
		return err
	}
	return nil
}
