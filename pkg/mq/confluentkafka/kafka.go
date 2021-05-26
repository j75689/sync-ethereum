package confluentkafka

import (
	"context"
	"encoding/json"
	"os"
	"strings"
	"sync"
	"sync-ethereum/pkg/mq"
	"sync-ethereum/pkg/util"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"
)

type _CallbackChannelMap struct {
	lock         sync.RWMutex
	callbackChan map[string]chan MsgData
}

func (m *_CallbackChannelMap) Get(key string) (chan MsgData, bool) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	c, ok := m.callbackChan[key]
	return c, ok
}

func (m *_CallbackChannelMap) Put(key string, value chan MsgData) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.callbackChan[key] = value
}

func (m *_CallbackChannelMap) Delete(key string) {
	m.lock.Lock()
	defer m.lock.Unlock()
	if c, ok := m.callbackChan[key]; ok {
		close(c)
	}
	delete(m.callbackChan, key)
}

func (m *_CallbackChannelMap) Close() {
	m.lock.Lock()
	defer m.lock.Unlock()
	for k, c := range m.callbackChan {
		close(c)
		delete(m.callbackChan, k)
	}
}

var _ mq.MQ = (*ConfluentKafka)(nil)

func NewConfluentKafka(option KafkaOption, logger zerolog.Logger) (*ConfluentKafka, error) {
	p, closeProducer, err := _NewProducer(option, logger)
	if err != nil {
		return nil, err
	}
	callbackChan := &_CallbackChannelMap{
		lock:         sync.RWMutex{},
		callbackChan: make(map[string]chan MsgData, 0),
	}
	c, closeConsumer, err := _NewConsumer(option, logger, callbackChan)
	if err != nil {
		return nil, err
	}
	return &ConfluentKafka{
		isRunning:     true,
		option:        option,
		producer:      p,
		consumer:      c,
		callbackChan:  callbackChan,
		closeConsumer: closeConsumer,
		closeProducer: closeProducer,
		subMiddleware: make([]func(key string, data []byte), 0),
	}, nil
}

func _NewConsumer(option KafkaOption, logger zerolog.Logger, callbackChan *_CallbackChannelMap) (*kafka.Consumer, func() error, error) {
	servers := strings.Join(option.Brokers, ",")
	if option.HeartbeatIntervalMs <= 0 {
		option.HeartbeatIntervalMs = 3000
	}
	if option.SessionTimeoutMs <= 0 {
		option.SessionTimeoutMs = 60000
	}
	if option.GroupID == "" {
		option.GroupID = uuid.New().String()
	}
	if option.OffsetsInitial == "" {
		option.OffsetsInitial = "earliest"
	}
	if option.GroupInstanceID == "" {
		gid := os.Getenv("HOSTNAME")
		if gid == "" {
			gid = uuid.New().String()
		}
		option.GroupInstanceID = gid
	}
	if option.RebalanceStrategy == "" {
		option.RebalanceStrategy = "cooperative-sticky"
	}
	if option.AutoCommitIntervalMs <= 0 {
		option.AutoCommitIntervalMs = 5000
	}
	if option.MaxPartitionFetchBytes <= 0 {
		option.MaxPartitionFetchBytes = 1048576
	}
	if option.FetchMaxBytes <= 0 {
		option.FetchMaxBytes = 1048576
	}
	if option.PollTimeoutMs <= 0 {
		option.PollTimeoutMs = 100
	}
	if option.SecurityProtoco == "" {
		option.SecurityProtoco = "plaintext"
	}
	if option.SASLMechanisms == "" {
		option.SASLMechanisms = "GSSAPI"
	}

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": servers,
		"group.id":          option.GroupID,
		// auto.offset.reset: latest, earliest, none
		// What to do when there is no initial offset in Kafka or if the current offset does not exist any more on the server (e.g. because that data has been deleted):
		// * earliest: automatically reset the offset to the earliest offset
		// * latest: automatically reset the offset to the latest offset
		// * none: throw exception to the consumer if no previous offset is found for the consumer's group
		"auto.offset.reset":             option.OffsetsInitial,
		"fetch.min.bytes":               1,
		"heartbeat.interval.ms":         option.HeartbeatIntervalMs,
		"session.timeout.ms":            option.SessionTimeoutMs,
		"enable.auto.commit":            option.EnableAutoCommit,
		"max.partition.fetch.bytes":     option.MaxPartitionFetchBytes,
		"fetch.max.bytes":               option.FetchMaxBytes,
		"group.instance.id":             option.GroupInstanceID,
		"max.poll.interval.ms":          300000,
		"partition.assignment.strategy": option.RebalanceStrategy,
		"auto.commit.interval.ms":       option.AutoCommitIntervalMs,
		"client.id":                     option.ClientID,
		"metadata.max.age.ms":           300000,
		"sasl.username":                 option.SASlUserName,
		"sasl.password":                 option.SASLPassword,
		"security.protocol":             option.SecurityProtoco,
		"sasl.mechanisms":               option.SASLMechanisms,
	})
	if err != nil {
		return nil, nil, err
	}
	isRunning := true
	go func() {
		defer callbackChan.Close()
		for isRunning {
			ev := c.Poll(option.PollTimeoutMs)
			if ev == nil {
				continue
			}
			switch e := ev.(type) {
			case *kafka.Message:
				receivedTopic := ""
				if e.TopicPartition.Topic != nil {
					receivedTopic = *e.TopicPartition.Topic
				}

				var msgData MsgData
				err := json.Unmarshal(e.Value, &msgData)
				if err != nil {
					logger.Error().Msgf("fail to unmarshal to internal msgData: %s", msgData)
					continue
				}
				msgData.ConsumeID = uuid.New().String()
				msgData.Commit = func() error {
					if option.EnableAutoCommit {
						return nil
					}
					_, err := c.CommitMessage(e)
					if err != nil {
						return errors.Wrap(err, "commit error")
					}
					return nil
				}
				logger.Debug().
					Str("request_id", msgData.RequestID).
					Str("consume_id", msgData.ConsumeID).
					Str("kafka.key", string(e.Key)).
					Str("kafka.topic", receivedTopic).
					Int64("kafka.offset", int64(e.TopicPartition.Offset)).
					Int32("kafka.partition", e.TopicPartition.Partition).
					Str("kafka.group", option.GroupID).
					Str("kafka.group.instance.id", option.GroupInstanceID).
					Str("kafka.msg.time", e.Timestamp.String()).
					Str("kafka.req", string(msgData.Data)).
					Msg("kafka access log")

				messageChan, ok := callbackChan.Get(receivedTopic)

				if ok && isRunning {
					messageChan <- msgData
				}

			case kafka.Error:
				// Errors should generally be considered
				// informational, the client will try to
				// automatically recover.
				// But in this example we choose to terminate
				// the application if all brokers are down.
				logger.Error().
					Str("kafka.group", option.GroupID).
					Str("kafka.group.instance.id", option.GroupInstanceID).
					Msgf("err: %v", e)

			case kafka.OffsetsCommitted:
				if e.Error != nil {
					logger.Error().Err(e.Error).Msg("offsets committed error")
				} else {
					logger.Info().Interface("offset_commit", e.Offsets).Msg("offsets committed")
				}
			default:
				logger.Info().
					Str("kafka.group", option.GroupID).
					Str("kafka.group.instance.id", option.GroupInstanceID).
					Msgf("confluent kafka ignored: %v", e)
			}
		}
	}()

	return c, func() error {
		isRunning = false
		if err := c.Close(); err != nil {
			return err
		}
		return nil
	}, nil
}

func _NewProducer(option KafkaOption, logger zerolog.Logger) (*kafka.Producer, func(), error) {
	servers := strings.Join(option.Brokers, ",")
	var acks string
	switch option.Acks {
	case 1:
		acks = "0"
	case 2:
		acks = "1"
	case 3:
		acks = "-1"
	default:
		acks = "1"
	}

	if option.CompressionType == "" {
		option.CompressionType = "none"
	}
	if option.BatchSize <= 0 {
		option.BatchSize = 1048576
	}
	if option.FlushWaitMs <= 0 {
		option.FlushWaitMs = 1000
	}
	if option.SecurityProtoco == "" {
		option.SecurityProtoco = "plaintext"
	}
	if option.SASLMechanisms == "" {
		option.SASLMechanisms = "GSSAPI"
	}

	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": servers,
		"acks":              acks,
		"compression.type":  option.CompressionType,
		"retries":           option.Retries,
		"batch.size":        option.BatchSize,
		"client.id":         option.ClientID,
		"sasl.username":     option.SASlUserName,
		"sasl.password":     option.SASLPassword,
		"security.protocol": option.SecurityProtoco,
		"sasl.mechanisms":   option.SASLMechanisms,
	})

	if err != nil {
		return nil, nil, err
	}

	closeProducer := func() {
		p.Close()
	}
	return p, closeProducer, nil
}

type ConfluentKafka struct {
	isRunning     bool
	option        KafkaOption
	producer      *kafka.Producer
	consumer      *kafka.Consumer
	callbackChan  *_CallbackChannelMap
	closeConsumer func() error
	closeProducer func()
	subMiddleware []func(key string, data []byte)
}

func (mq *ConfluentKafka) Publish(topic, key string, data []byte) error {
	if !mq.isRunning {
		return nil
	}
	msgData := MsgData{
		RequestID: key,
		Data:      data,
	}

	b, err := json.Marshal(msgData)
	if err != nil {
		return errors.New("fail to marshal internal MsgData, err: " + err.Error())
	}

	kafkaMsg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Key:   []byte(key),
		Value: b,
	}

	mq.producer.ProduceChannel() <- kafkaMsg
	return nil
}

func (mq *ConfluentKafka) Subscribe(ctx context.Context, workerSize int, topic string, process func(key string, data []byte) (bool, error), errCallBack ...func(string, error)) error {
	if !mq.isRunning {
		return nil
	}
	messageChan := make(chan MsgData, 0)
	mq.callbackChan.Put(topic, messageChan)

	err := mq.consumer.Subscribe(topic, nil)
	if err != nil {
		return err
	}
	return mq._StartSubscribeWorker(ctx, workerSize, messageChan, process, errCallBack...)
}

func (mq *ConfluentKafka) _StartSubscribeWorker(ctx context.Context, workerSize int, messageChan <-chan MsgData, process func(key string, data []byte) (bool, error), errCallBack ...func(string, error)) error {
	errGroup := errgroup.Group{}
	for i := 0; i < workerSize; i++ {
		errGroup.Go(func() (err error) {
			defer func() {
				if recoverErr := util.ConvertRecoverToError(recover()); recoverErr != nil {
					err = recoverErr
				}
			}()

			for mq.isRunning {
				select {
				case m := <-messageChan:
					if !mq.isRunning {
						return
					}
					for _, mid := range mq.subMiddleware {
						mid(m.RequestID, m.Data)
					}

					// recover panic
					f := func(key string, data []byte) (ack bool, err error) {
						defer func() {
							if recoverErr := util.ConvertRecoverToError(recover()); recoverErr != nil {
								err = recoverErr
								ack = true
							}
						}()
						return process(key, data)
					}

					isAck, err := f(m.RequestID, m.Data)
					if err != nil {
						for _, cb := range errCallBack {
							cb(m.RequestID, err)
						}
					}
					if isAck && mq.isRunning {
						err := m.Commit()
						if err != nil {
							for _, cb := range errCallBack {
								cb(m.RequestID, err)
							}
						}
					}
				case <-ctx.Done():
					return
				}
			}
			return
		})
	}

	return errGroup.Wait()
}

func (mq *ConfluentKafka) SubscriberMiddleware(middleware ...func(key string, data []byte)) {
	mq.subMiddleware = append(mq.subMiddleware, middleware...)
}

func (mq *ConfluentKafka) Close() error {
	mq.isRunning = false
	if mq.closeConsumer != nil {
		err := mq.closeConsumer()
		if err != nil {
			return err
		}
	}
	if mq.closeProducer != nil {
		mq.closeProducer()
	}
	return nil
}
