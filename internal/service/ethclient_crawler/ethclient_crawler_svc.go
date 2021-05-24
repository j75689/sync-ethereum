package ethclient_crawler

import (
	"context"
	"math/big"
	"sync-ethereum/internal/config"
	"sync-ethereum/internal/service"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

var _ service.CrawlerService = (*EthClientCrawlerService)(nil)

func NewEthClientCrawlerService(config config.Config) service.CrawlerService {
	return &EthClientCrawlerService{
		clientPool: _NewClientPool(config.EthClient.DialTimeout, config.EthClient.URL, config.EthClient.MaxClientConn),
	}
}

type EthClientCrawlerService struct {
	clientPool *_ClientPool
}

func (svc *EthClientCrawlerService) GetBlockNumber(ctx context.Context) (*big.Int, error) {
	client, err := svc.clientPool.Get()
	if err != nil {
		return nil, err
	}

	number, err := client.BlockNumber(ctx)
	if err != nil {
		return nil, err
	}
	return big.NewInt(int64(number)), nil
}

func (svc *EthClientCrawlerService) GetBlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error) {
	client, err := svc.clientPool.Get()
	if err != nil {
		return nil, err
	}

	return client.BlockByNumber(ctx, number)
}

func (svc *EthClientCrawlerService) GetTransactionByHash(ctx context.Context, hash common.Hash) (tx *types.Transaction, isPending bool, err error) {
	client, err := svc.clientPool.Get()
	if err != nil {
		return nil, false, err
	}

	return client.TransactionByHash(ctx, hash)
}

func (svc *EthClientCrawlerService) GetTransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	client, err := svc.clientPool.Get()
	if err != nil {
		return nil, err
	}

	return client.TransactionReceipt(ctx, txHash)
}

func (svc *EthClientCrawlerService) Close() {
	svc.clientPool.Close()
}
