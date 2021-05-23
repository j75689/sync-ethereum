package storage

import (
	"context"
	"sync-ethereum/internal/model"
	"sync-ethereum/internal/repository"
	"sync-ethereum/internal/service"
)

var _ service.StorageService = (*StorageService)(nil)

func NewStorageService(repo repository.StorageRepository) service.StorageService {
	return &StorageService{
		repo: repo,
	}
}

type StorageService struct {
	repo repository.StorageRepository
}

func (svc *StorageService) GetCurrentBlockNumber(ctx context.Context) (model.GormBigInt, error) {
	currentBlockNumber, err := svc.repo.GetCurrentBlockNumber(ctx)
	if err != nil {
		return model.GormBigInt{}, err
	}
	return currentBlockNumber.BlockNumber, nil
}

func (svc *StorageService) UpdateCurrentBlockNumber(ctx context.Context, blockNumber model.GormBigInt, onlineBlockNumber model.GormBigInt) error {
	return svc.repo.UpdateCurrentBlockNumber(ctx, &model.CurrentBlockNumber{BlockNumber: blockNumber, OnlineBlockNumber: onlineBlockNumber})
}

func (svc *StorageService) GetBlock(ctx context.Context, filter model.Block) (model.Block, error) {
	return svc.repo.GetBlock(ctx, filter, model.Block{}.Preload)
}

func (svc *StorageService) ListBlock(ctx context.Context, filter model.Block, pagination model.Pagination, sorting model.Sorting) ([]model.Block, error) {
	return svc.repo.ListBlock(ctx, filter, pagination.LimitAndOffset, sorting.Sort, model.Block{}.Preload)
}

func (svc *StorageService) CreateBlock(ctx context.Context, block *model.Block) error {
	return svc.repo.CreateBlock(ctx, block, model.Block{}.OnConflict)
}

func (svc *StorageService) UpdateBlock(ctx context.Context, filter model.Block, block *model.Block) error {
	return svc.repo.UpdateBlock(ctx, filter, block)
}

func (svc *StorageService) GetTransaction(ctx context.Context, filter model.Transaction) (model.Transaction, error) {
	return svc.repo.GetTransaction(ctx, filter, model.Transaction{}.Preload)
}

func (svc *StorageService) Close() error {
	return svc.repo.Close()
}
