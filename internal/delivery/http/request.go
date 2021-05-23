package http

import "sync-ethereum/internal/model"

type GetBlocksResponse struct {
	Blocks []Block `json:"block"`
}

type Block struct {
	BlockNumber model.GormBigInt `json:"block_num"`
	BlockHash   string           `json:"block_hash"`
	BlockTime   uint64           `json:"block_time"`
	ParentHash  string           `json:"parent_hash"`
	IsStable    bool             `json:"is_stable"`
}

type GetBlockResponse struct {
	BlockNumber  model.GormBigInt `json:"block_num"`
	BlockHash    string           `json:"block_hash"`
	BlockTime    uint64           `json:"block_time"`
	ParentHash   string           `json:"parent_hash"`
	IsStable     bool             `json:"is_stable"`
	Transactions []string         `json:"transactions"`
}

type GetTransactionResponse struct {
	TXHash string           `json:"tx_hash"`
	From   string           `json:"from"`
	To     string           `json:"to"`
	Nonce  uint64           `json:"nonce"`
	Data   model.BlockData  `json:"data"`
	Value  model.GormBigInt `json:"value"`
	Logs   []TransactionLog `json:"logs"`
}

type TransactionLog struct {
	Index uint64          `json:"index"`
	Data  model.BlockData `json:"data"`
}
