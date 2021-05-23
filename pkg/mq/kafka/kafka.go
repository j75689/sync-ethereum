package kafka

import (
	"context"
	"errors"
	"sync-ethereum/pkg/mq"
	"sync-ethereum/pkg/util"

	"github.com/Shopify/sarama"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
	"github.com/ThreeDotsLabs/watermill/message"
	"golang.org/x/sync/errgroup"
)

type KafkaOption struct {
	Brokers        []string
	ConsumerGroup  string
	OffsetsInitial int64
	FetchDefault   int32
	RequiredAcks   int16
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
	saramaSubscriberConfig.Consumer.Fetch.Default = option.FetchDefault
	saramaSubscriberConfig.Producer.RequiredAcks = sarama.RequiredAcks(option.RequiredAcks)
	subscriber, err := kafka.NewSubscriber(
		kafka.SubscriberConfig{
			Brokers:               option.Brokers,
			Unmarshaler:           kafka.DefaultMarshaler{},
			OverwriteSaramaConfig: saramaSubscriberConfig,
			ConsumerGroup:         option.ConsumerGroup,
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

func (mq *KafkaMQ) Subscribe(ctx context.Context, workerSize int, topic string, process func(key string, data []byte) (bool, error), errCallBack ...func(string, error)) error {
	if process == nil {
		return errors.New("process is nil function")
	}

	message, err := mq.subscriber.Subscribe(ctx, topic)
	if err != nil {
		return err
	}

	return mq._StartSubscribeWorker(ctx, workerSize, message, process, errCallBack...)
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

func (mq *KafkaMQ) _StartSubscribeWorker(ctx context.Context, workerSize int, messageChan <-chan *message.Message,
	process func(key string, data []byte) (bool, error), errCallBack ...func(string, error)) error {

	// non supporte worker
	errGroup := errgroup.Group{}
	errGroup.Go(func() (err error) {
		defer func() {
			if recoverErr := util.ConvertRecoverToError(recover()); recoverErr != nil {
				err = recoverErr
			}
		}()

		for {
			select {
			case m := <-messageChan:
				if m == nil {
					continue
				}

				for _, mid := range mq.subMiddleware {
					mid(m.UUID, m.Payload)
				}

				// recover panic
				f := func(key string, data []byte) (ack bool, err error) {
					defer func() {
						if recoverErr := util.ConvertRecoverToError(recover()); recoverErr != nil {
							err = recoverErr
							ack = true
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
			case <-ctx.Done():
				return
			}
		}
	})

	return errGroup.Wait()
}
