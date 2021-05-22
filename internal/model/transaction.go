package model

import (
	"time"
)

type Transaction struct {
	TXHash      string            `json:"tx_hash" gorm:"type:varchar(128);column:tx_hash;primaryKey;autoIncrement:false"`
	BlockNumber GormBigInt        `json:"block_num" gorm:"type:bigint;column:block_num;index"`
	From        string            `json:"from" gorm:"type:varchar(128)"`
	To          string            `json:"to" gorm:"type:varchar(128)"`
	Nonce       uint64            `json:"nonce"`
	Data        BlockData         `json:"data"`
	Value       GormBigInt        `json:"value" gorm:"type:bigint"`
	Logs        []*TransactionLog `json:"logs" gorm:"foreignKey:TXHash;references:TXHash"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	DeletedAt   *time.Time        `json:"deleted_at" gorm:"index"`
}

type TransactionLog struct {
	ID        int64      `json:"id" gorm:"primaryKey"`
	TXHash    string     `json:"tx_hash" gorm:"type:varchar(128)column:tx_hash;index"`
	Index     uint64     `json:"index"`
	Data      BlockData  `json:"data"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" gorm:"index"`
}
