package mq

import (
	"context"
)

type MQ interface {
	Publish(topic, key string, data []byte) error
	Subscribe(ctx context.Context, topic string, process func(key string, data []byte) (bool, error), errCallBack ...func(string, error)) error
	SubscriberMiddleware(middleware ...func(key string, data []byte))
	Close() error
}
