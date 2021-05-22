package gorm

import (
	"context"
	"sync-ethereum/internal/model"
	"sync-ethereum/internal/repository"
	"sync-ethereum/internal/repository/gorm/migration"

	"github.com/go-gormigrate/gormigrate/v2"

	"gorm.io/gorm"
)

var _ repository.StorageRepository = (*StorageRepository)(nil)

func NewStorageRepository(db *gorm.DB) repository.StorageRepository {
	return &StorageRepository{
		migration: gormigrate.New(db, gormigrate.DefaultOptions, migration.Migrations),
		db:        db,
	}
}

type StorageRepository struct {
	migration *gormigrate.Gormigrate
	db        *gorm.DB
}

func (repo *StorageRepository) MigrateUp() error {
	return repo.migration.Migrate()
}

func (repo *StorageRepository) MigrateDown() error {
	for _, m := range migration.Migrations {
		if err := repo.migration.RollbackMigration(m); err != nil {
			return err
		}
	}
	return nil
}

func (repo *StorageRepository) MigrateUpTo(version string) error {
	return repo.migration.MigrateTo(version)
}

func (repo *StorageRepository) MigrateDownTo(version string) error {
	return repo.migration.RollbackTo(version)
}

func (repo *StorageRepository) GetBlock(ctx context.Context, filter model.Block, scope ...func(*gorm.DB) *gorm.DB) (model.Block, error) {
	block := model.Block{}
	tx := repo.db.WithContext(ctx).Scopes(scope...).Where(filter).First(&block)
	return block, tx.Error
}

func (repo *StorageRepository) ListBlock(ctx context.Context, filter model.Block, scope ...func(*gorm.DB) *gorm.DB) ([]model.Block, error) {
	blocks := []model.Block{}
	tx := repo.db.WithContext(ctx).Scopes(scope...).Model(model.Block{}).Where(filter).Find(&blocks)
	return blocks, tx.Error
}

func (repo *StorageRepository) CreateBlock(ctx context.Context, block *model.Block, scope ...func(*gorm.DB) *gorm.DB) error {
	return repo.db.WithContext(ctx).Scopes(scope...).Create(block).Error
}

func (repo *StorageRepository) UpdateBlock(ctx context.Context, filter model.Block, block *model.Block, scope ...func(*gorm.DB) *gorm.DB) error {
	return repo.db.WithContext(ctx).Scopes(scope...).Where(filter).Updates(block).Error
}

func (repo *StorageRepository) GetTransaction(ctx context.Context, filter model.Transaction, scope ...func(*gorm.DB) *gorm.DB) (model.Transaction, error) {
	transaction := model.Transaction{}
	tx := repo.db.WithContext(ctx).Scopes(scope...).Where(filter).First(&transaction)
	return transaction, tx.Error
}

func (repo *StorageRepository) Close() error {
	db, err := repo.db.DB()
	if err != nil {
		return err
	}

	return db.Close()
}
