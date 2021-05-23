package wireset

import (
	"strings"
	"sync-ethereum/internal/config"
	"sync-ethereum/pkg/mq"
	"sync-ethereum/pkg/mq/confluentkafka"
	"sync-ethereum/pkg/mq/kafka"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

func InitMQ(config config.Config, log zerolog.Logger) (queue mq.MQ, err error) {
	switch strings.ToLower(config.MQ.Driver) {
	case "kafka":
		queue, err = kafka.NewKafkaMQ(kafka.KafkaOption{
			Brokers:        config.MQ.KafkaOption.Brokers,
			ConsumerGroup:  config.MQ.KafkaOption.ConsumerGroup,
			OffsetsInitial: config.MQ.KafkaOption.OffsetsInitial,
			FetchDefault:   config.MQ.KafkaOption.FetchDefault,
			RequiredAcks:   config.MQ.KafkaOption.RequiredAcks,
			LoggerAdapter:  mq.WrapWatermillLogger(log),
		})
	case "confluentkafka":
		queue, err = confluentkafka.NewConfluentKafka(confluentkafka.KafkaOption{
			Brokers:                config.MQ.ConfluentKafkaOption.Brokers,
			ClientID:               config.MQ.ConfluentKafkaOption.ClientID,
			ConsumerGroup:          config.MQ.ConfluentKafkaOption.ConsumerGroup,
			GroupID:                config.MQ.ConfluentKafkaOption.GroupID,
			OffsetsInitial:         config.MQ.ConfluentKafkaOption.OffsetsInitial,
			RebalanceStrategy:      config.MQ.ConfluentKafkaOption.RebalanceStrategy,
			GroupInstanceID:        config.MQ.ConfluentKafkaOption.GroupInstanceID,
			HeartbeatIntervalMs:    config.MQ.ConfluentKafkaOption.HeartbeatIntervalMs,
			SessionTimeoutMs:       config.MQ.ConfluentKafkaOption.SessionTimeoutMs,
			AutoCommitIntervalMs:   config.MQ.ConfluentKafkaOption.AutoCommitIntervalMs,
			EnableAutoCommit:       config.MQ.ConfluentKafkaOption.EnableAutoCommit,
			Acks:                   config.MQ.ConfluentKafkaOption.Acks,
			CompressionType:        config.MQ.ConfluentKafkaOption.CompressionType,
			Retries:                config.MQ.ConfluentKafkaOption.Retries,
			BatchSize:              config.MQ.ConfluentKafkaOption.BatchSize,
			FlushWaitMs:            config.MQ.ConfluentKafkaOption.FlushWaitMs,
			FetchMaxBytes:          config.MQ.ConfluentKafkaOption.FetchMaxBytes,
			MaxPartitionFetchBytes: config.MQ.ConfluentKafkaOption.MaxPartitionFetchBytes,
			PollTimeoutMs:          config.MQ.ConfluentKafkaOption.PollTimeoutMs,
		}, log)
	default:
		err = errors.New("no supported driver [" + config.MQ.Driver + "]")
	}
	if queue != nil && err == nil {
		queue.SubscriberMiddleware(func(key string, data []byte) {
			log.Info().Str("message_key", key).Bytes("message", data).Send()
		})
	}
	return
}
