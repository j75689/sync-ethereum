package model

type CurrentBlockNumber struct {
	ID          int64      `json:"id" gorm:"primaryKey"`
	BlockNumber GormBigInt `json:"block_num" gorm:"type:bigint;column:block_num"`
}
