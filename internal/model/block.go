package model

import (
	"time"
)

type Block struct {
	BlockNumber GormBigInt     `json:"block_num" gorm:"type:bigint;column:block_num;primaryKey;autoIncrement:false"`
	BlockHash   string         `json:"block_hash" gorm:"type:varchar(128);column:block_hash;uniqueIndex"`
	BlockTime   uint64         `json:"block_time"`
	ParentHash  string         `json:"parent_hash" gorm:"type:varchar(128);column:parent_hash;uniqueIndex"`
	IsStable    bool           `json:"is_stable"`
	Transaction []*Transaction `gorm:"foreignKey:BlockNumber;references:BlockNumber"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   *time.Time     `json:"deleted_at" gorm:"index"`
}
