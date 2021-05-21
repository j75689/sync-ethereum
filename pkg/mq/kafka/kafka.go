package kafka

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"sync-ethereum/pkg/mq"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
	"github.com/ThreeDotsLabs/watermill/message"
)

type KafkaOption struct {
	Brokers        []string
	ConsumerGroup  string
	OffsetsInitial int64
	LoggerAdapter  watermill.LoggerAdapter
}

var _ mq.MQ = (*KafkaMQ)(nil)

func NewKafkaMQ(option KafkaOption) (*KafkaMQ, error) {
	publisher, err := kafka.NewPublisher(
		kafka.PublisherConfig{
			Brokers:   option.Brokers,
			Marshaler: kafka.DefaultMarshaler{},
		}, option.LoggerAdapter,
	)
	if err != nil {
		return nil, err
	}

	saramaSubscriberConfig := kafka.DefaultSaramaSubscriberConfig()
	saramaSubscriberConfig.Consumer.Offsets.Initial = option.OffsetsInitial
	subscriber, err := kafka.NewSubscriber(
		kafka.SubscriberConfig{
			Brokers:       option.Brokers,
			Unmarshaler:   kafka.DefaultMarshaler{},
			ConsumerGroup: option.ConsumerGroup,
		}, option.LoggerAdapter,
	)
	if err != nil {
		return nil, err
	}

	return &KafkaMQ{
		publisher:     publisher,
		subscriber:    subscriber,
		subMiddleware: make([]func(key string, data []byte), 0),
	}, nil
}

type KafkaMQ struct {
	publisher     message.Publisher
	subscriber    message.Subscriber
	subMiddleware []func(key string, data []byte)
}

func (mq *KafkaMQ) Publish(topic, key string, data []byte) error {
	if len(key) == 0 {
		key = watermill.NewUUID()
	}
	return mq.publisher.Publish(topic, message.NewMessage(key, message.Payload(data)))
}

func (mq *KafkaMQ) Subscribe(ctx context.Context, topic string, process func(key string, data []byte) (bool, error), errCallBack ...func(string, error)) error {
	if process == nil {
		return errors.New("process is nil function")
	}

	message, err := mq.subscriber.Subscribe(ctx, topic)
	if err != nil {
		return err
	}

	for m := range message {
		for _, mid := range mq.subMiddleware {
			mid(m.UUID, m.Payload)
		}

		// recover panic
		f := func(key string, data []byte) (ack bool, err error) {
			defer func() {
				if r := recover(); r != nil {
					var msg string
					for i := 2; ; i++ {
						_, file, line, ok := runtime.Caller(i)
						if !ok {
							break
						}
						msg += fmt.Sprintf("%s:%d\n", file, line)
					}
					ack = true
					err = errors.New(msg)
				}
			}()
			return process(m.UUID, m.Payload)
		}

		isAck, err := f(m.UUID, m.Payload)
		if err != nil {
			for _, cb := range errCallBack {
				cb(m.UUID, err)
			}
		}
		if isAck {
			m.Ack()
		} else {
			m.Nack()
		}
	}

	return nil
}

func (mq *KafkaMQ) SubscriberMiddleware(middleware ...func(key string, data []byte)) {
	mq.subMiddleware = append(mq.subMiddleware, middleware...)
}

func (mq *KafkaMQ) Close() error {
	if err := mq.publisher.Close(); err != nil {
		return err
	}
	return mq.subscriber.Close()
}
