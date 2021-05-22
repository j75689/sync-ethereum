package ethclient_crawler

import (
	"context"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
)

func _NewClientPool(dialTimeout time.Duration, url string) *_ClientPool {
	clientPool := &_ClientPool{
		clients: make([]*ethclient.Client, 0),
	}

	clientPool.clientPool = sync.Pool{
		New: func() interface{} {
			ctx, cancel := context.WithTimeout(context.Background(), dialTimeout)
			defer cancel()
			client, err := ethclient.DialContext(ctx, url)
			if err != nil {
				return err
			}
			clientPool.clients = append(clientPool.clients, client)
			return client
		},
	}

	return clientPool
}

type _ClientPool struct {
	clients    []*ethclient.Client
	clientPool sync.Pool
}

func (pool *_ClientPool) Get() (*ethclient.Client, error) {
	v := pool.clientPool.Get()
	client, ok := v.(*ethclient.Client)
	if !ok {
		return nil, v.(error)
	}

	return client, nil
}

func (pool *_ClientPool) Put(client *ethclient.Client) {
	pool.clientPool.Put(client)
}

func (pool *_ClientPool) Len() int {
	return len(pool.clients)
}

func (pool *_ClientPool) Close() {
	for _, client := range pool.clients {
		if client != nil {
			client.Close()
		}
	}
}
