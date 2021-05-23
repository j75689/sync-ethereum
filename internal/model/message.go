package model

type CrawlerMessage struct {
	IsStable    bool       `json:"is_stable"`
	BlockNumber GormBigInt `json:"block_number"`
}
