package service

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type CrawlerService interface {
	GetBlockNumber(ctx context.Context) (*big.Int, error)
	GetBlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error)
	GetTransactionByHash(ctx context.Context, hash common.Hash) (tx *types.Transaction, isPending bool, err error)
	GetTransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error)
	Close()
}
