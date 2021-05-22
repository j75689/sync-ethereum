package repository

import (
	"context"
	"sync-ethereum/internal/model"

	"gorm.io/gorm"
)

type StorageRepository interface {
	MigrateUp() error
	MigrateDown() error
	MigrateUpTo(version string) error
	MigrateDownTo(version string) error
	GetCurrentBlockNumber(ctx context.Context, scope ...func(*gorm.DB) *gorm.DB) (model.CurrentBlockNumber, error)
	UpdateCurrentBlockNumber(ctx context.Context, blockNumber *model.CurrentBlockNumber, scope ...func(*gorm.DB) *gorm.DB) error
	GetBlock(ctx context.Context, filter model.Block, scope ...func(*gorm.DB) *gorm.DB) (model.Block, error)
	ListBlock(ctx context.Context, filter model.Block, scope ...func(*gorm.DB) *gorm.DB) ([]model.Block, error)
	CreateBlock(ctx context.Context, block *model.Block, scope ...func(*gorm.DB) *gorm.DB) error
	UpdateBlock(ctx context.Context, filter model.Block, block *model.Block, scope ...func(*gorm.DB) *gorm.DB) error
	GetTransaction(ctx context.Context, filter model.Transaction, scope ...func(*gorm.DB) *gorm.DB) (model.Transaction, error)
	Close() error
}
