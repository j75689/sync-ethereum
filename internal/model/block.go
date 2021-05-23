package model

import (
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Block struct {
	BlockNumber GormBigInt     `json:"block_num" gorm:"type:varchar(32);column:block_num;primaryKey;autoIncrement:false"`
	BlockHash   string         `json:"block_hash" gorm:"type:varchar(128);column:block_hash;uniqueIndex:idx_block_parent_hash"`
	BlockTime   uint64         `json:"block_time"`
	ParentHash  string         `json:"parent_hash" gorm:"type:varchar(128);column:parent_hash;uniqueIndex:idx_block_parent_hash"`
	IsStable    bool           `json:"is_stable"`
	Transaction []*Transaction `gorm:"foreignKey:BlockNumber;references:BlockNumber"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   *time.Time     `json:"deleted_at" gorm:"index"`
}

func (block Block) Preload(db *gorm.DB) *gorm.DB {
	return db.Preload("Transaction")
}

func (block Block) OnConflict(db *gorm.DB) *gorm.DB {
	return db.Clauses(clause.OnConflict{
		UpdateAll: true,
	})
}
