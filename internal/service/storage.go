package service

import (
	"context"
	"sync-ethereum/internal/model"
)

type StorageService interface {
	GetBlock(ctx context.Context, filter model.Block) (model.Block, error)
	ListBlock(ctx context.Context, filter model.Block, pagination model.Pagination, sorting model.Sorting) ([]model.Block, error)
	CreateBlock(ctx context.Context, block *model.Block) error
	UpdateBlock(ctx context.Context, filter model.Block, block *model.Block) error
	GetTransaction(ctx context.Context, filter model.Transaction) (model.Transaction, error)
	Close() error
}
