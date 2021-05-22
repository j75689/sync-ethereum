package wireset

import (
	"strings"
	"sync-ethereum/internal/config"
	"sync-ethereum/pkg/mq"
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
			LoggerAdapter:  mq.WrapWatermillLogger(log),
		})
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
