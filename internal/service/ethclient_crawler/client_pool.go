package ethclient_crawler

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
)

func _NewClient(dialTimeout time.Duration, url string, clientNum *int32) *_Client {
	return &_Client{
		lock:        sync.Mutex{},
		dialTimeout: dialTimeout,
		url:         url,
		clientNum:   clientNum,
	}
}

type _Client struct {
	lock        sync.Mutex
	ethclient   *ethclient.Client
	dialTimeout time.Duration
	url         string
	clientNum   *int32
}

func (c *_Client) Get() (*ethclient.Client, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.ethclient != nil {
		return c.ethclient, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), c.dialTimeout)
	defer cancel()
	client, err := ethclient.DialContext(ctx, c.url)
	if err != nil {
		return nil, err
	}
	defer func() {
		c.ethclient = client
		atomic.AddInt32(c.clientNum, 1)
	}()
	return client, err
}

func (c *_Client) Close() {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.ethclient != nil {
		c.ethclient.Close()
		c.ethclient = nil
	}
}

func _NewClientPool(dialTimeout time.Duration, url string, maxClientConn int) *_ClientPool {
	if maxClientConn <= 0 {
		maxClientConn = 1 // default
	}
	var clientNum int32

	clients := make([]*_Client, maxClientConn)
	for i := 0; i < maxClientConn; i++ {
		clients[i] = _NewClient(dialTimeout, url, &clientNum)
	}
	clientPool := &_ClientPool{
		maxClientNum: int32(maxClientConn),
		clientNum:    &clientNum,
		clients:      clients,
	}
	return clientPool
}

type _ClientPool struct {
	roundRobin   int32
	clientNum    *int32
	maxClientNum int32
	clients      []*_Client
}

func (pool *_ClientPool) Get() (*ethclient.Client, error) {
	current := atomic.LoadInt32(&pool.roundRobin)
	if current >= pool.maxClientNum {
		atomic.StoreInt32(&pool.roundRobin, 0)
		current = 0
	}

	defer func() {
		atomic.AddInt32(&pool.roundRobin, 1)
	}()
	return pool.clients[current].Get()
}

func (pool *_ClientPool) Len() int {
	return int(atomic.LoadInt32(pool.clientNum))
}

func (pool *_ClientPool) Close() {
	for _, client := range pool.clients {
		if client != nil {
			client.Close()
		}
	}
}
