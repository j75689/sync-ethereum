package model

type CurrentBlockNumber struct {
	ID                int64      `json:"id" gorm:"primaryKey"`
	BlockNumber       GormBigInt `json:"block_num" gorm:"type:varchar(32);column:block_num"`
	OnlineBlockNumber GormBigInt `json:"online_block_num" gorm:"type:varchar(32);column:online_block_num"`
}
