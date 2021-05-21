package wireset

import (
	"context"
	"sync-ethereum/internal/config"

	"github.com/ethereum/go-ethereum/ethclient"
)

func InitEthClient(config config.Config) (*ethclient.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.EthClient.DialTimeout)
	defer cancel()
	return ethclient.DialContext(ctx, config.EthClient.URL)
}
