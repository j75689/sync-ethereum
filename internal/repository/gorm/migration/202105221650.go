package migration

import (
	"math/big"
	"sync-ethereum/internal/model"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

var v202105221650 = &gormigrate.Migration{
	ID: "202105221650",
	Migrate: func(tx *gorm.DB) error {
		if err := tx.AutoMigrate(&model.Block{}); err != nil {
			return err
		}
		if err := tx.AutoMigrate(&model.Transaction{}); err != nil {
			return err
		}
		if err := tx.AutoMigrate(&model.TransactionLog{}); err != nil {
			return err
		}
		if err := tx.AutoMigrate(&model.CurrentBlockNumber{}); err != nil {
			return err
		}
		blockNumber := big.NewInt(0)
		if err := tx.Create(&model.CurrentBlockNumber{
			BlockNumber: model.GormBigInt(*blockNumber),
		}).Error; err != nil {
			return err
		}
		return nil
	},
	Rollback: func(tx *gorm.DB) error {
		if err := tx.Migrator().DropTable(&model.Block{}); err != nil {
			return err
		}
		if err := tx.Migrator().DropTable(&model.Transaction{}); err != nil {
			return err
		}
		if err := tx.Migrator().DropTable(&model.TransactionLog{}); err != nil {
			return err
		}
		if err := tx.Migrator().DropTable(&model.CurrentBlockNumber{}); err != nil {
			return err
		}
		return nil
	},
}
