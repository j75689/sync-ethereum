package config

import (
	"os"
	"strings"
	"sync-ethereum/pkg/logger"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	APPID          string               `mapstructure:"app_id"`
	Release        bool                 `mapstructure:"release"`
	Logger         LoggerConfig         `mapstructure:"logger"`
	HTTP           HTTPConfig           `mapstructure:"http"`
	DataBase       DataBaseConfig       `mapstructure:"database"`
	MQ             MQConfig             `mapstructure:"mq"`
	EthClient      EthClientConfig      `mapstructure:"eth_client"`
	Scheduler      SchedulerConfig      `mapstructure:"scheduler"`
	Crawler        CrawlerConfig        `mapstructure:"crawler"`
	DatabaseWriter DatabaseWriterConfig `mapstructure:"database_writer"`
}

type LoggerConfig struct {
	Level  string           `mapstructure:"level"`
	Format logger.LogFormat `mapstructure:"format"`
}

type HTTPConfig struct {
	Port uint16 `mapstructure:"port"`
}

type DataBaseConfig struct {
	Driver         string        `mapstructure:"driver"`
	Host           string        `mapstructure:"host"`
	Port           uint          `mapstructure:"port"`
	Database       string        `mapstructure:"database"`
	InstanceName   string        `mapstructure:"instance_name"`
	User           string        `mapstructure:"user"`
	Password       string        `mapstructure:"password"`
	ConnectTimeout string        `mapstructure:"connect_timeout"`
	ReadTimeout    string        `mapstructure:"read_timeout"`
	WriteTimeout   string        `mapstructure:"write_timeout"`
	DialTimeout    time.Duration `mapstructure:"dial_timeout"`
	MaxLifetime    time.Duration `mapstructure:"max_lifetime"`
	MaxIdleTime    time.Duration `mapstructure:"max_idletime"`
	MaxIdleConn    int           `mapstructure:"max_idle_conn"`
	MaxOpenConn    int           `mapstructure:"max_open_conn"`
	SSLMode        bool          `mapstructure:"ssl_mode"`
}

type MQConfig struct {
	Driver               string                     `mapstructure:"driver"`
	KafkaOption          KafkaOptionConfig          `mapstructure:"kafka_option"`
	ConfluentKafkaOption ConfluentKafkaOptionConfig `mapstructure:"confluentkafka_option"`
}

type KafkaOptionConfig struct {
	Brokers        []string `mapstructure:"brokers"`
	ConsumerGroup  string   `mapstructure:"consumer_group"`
	OffsetsInitial int64    `mapstructure:"offsets_initial"`
	FetchDefault   int32    `mapstructure:"fetch_default"`
	RequiredAcks   int16    `mapstructure:"required_acks"`
}

type ConfluentKafkaOptionConfig struct {
	Brokers       []string `mapstructure:"brokers"`
	ClientID      string   `mapstructure:"client_id"`
	ConsumerGroup string   `mapstructure:"consumer_group"`
	// All clients sharing the same group.id belong to the same group.
	GroupID string `mapstructure:"group_id"`
	// earliest / latest
	OffsetsInitial string `mapstructure:"offsets_initial"`
	// range / roundrobin / cooperative-sticky; default cooperative-sticky
	RebalanceStrategy    string `mapstructure:"rebalance_strategy"`
	GroupInstanceID      string `mapstructure:"group_instance_id"`
	HeartbeatIntervalMs  int    `mapstructure:"heartbeat_interval_ms"`
	SessionTimeoutMs     int    `mapstructure:"session_timeout_ms"`
	AutoCommitIntervalMs int    `mapstructure:"auto_commit_interval_ms"`
	EnableAutoCommit     bool   `mapstructure:"enable_auto_commit"`

	// ================ producer related config ================
	// valid value: 1、2、3;
	// 1:
	// The producer will not wait for any acknowledgment from the server at all.
	// The record will be immediately added to the socket buffer and considered sent.
	// 2:
	// This will mean the leader will write the record to its local log but will respond without awaiting full acknowledgement from all followers.
	// 3:
	// This means the leader will wait for the full set of in-sync replicas (ISR) to acknowledge the record.
	Acks int `mapstructure:"acks"`
	// none, gzip, snappy, lz4, zstd
	CompressionType string `mapstructure:"compression_type"`
	Retries         int    `mapstructure:"retries"`
	BatchSize       int    `mapstructure:"batch_size"`
	FlushWaitMs     int    `mapstructure:"flush_wait_ms"`

	// ================ consumer related config ================
	FetchMaxBytes          int `mapstructure:"fetch_max_bytes"`
	MaxPartitionFetchBytes int `mapstructure:"max_partition_fetch_bytes"`
	PollTimeoutMs          int `mapstructure:"poll_timeout_ms"`

	// ================ auth ================
	SASlUserName    string `mapstructure:"sasl_username"`
	SASLPassword    string `mapstructure:"sasl_password"`
	SASLMechanisms  string `mapstructure:"sasl_mechanisms"`
	SecurityProtoco string `mapstructure:"security_protoco"`
}

type EthClientConfig struct {
	URL           string        `mapstructure:"url"`
	DialTimeout   time.Duration `mapstructure:"dial_timeout"`
	MaxClientConn int           `mapstructure:"max_client_conn"`
}

type SchedulerConfig struct {
	UnstableNumber int        `mapstructure:"unstable_num"`
	StartAt        int64      `mapstructure:"start_at"`
	Sync           SyncConfig `mapstructure:"sync"`
	BatchLimit     int64      `mapstructure:"batch_limit"`
}
type SyncConfig struct {
	Interval time.Duration `mapstructure:"interval"`
}
type CrawlerConfig struct {
	Topic    string        `mapstructure:"topic"`
	PoolSize int           `mapstructure:"pool_size"`
	Timeout  time.Duration `mapstructure:"timeout"`
}

type DatabaseWriterConfig struct {
	Topic    string        `mapstructure:"topic"`
	PoolSize int           `mapstructure:"pool_size"`
	Timeout  time.Duration `mapstructure:"timeout"`
}

func NewConfig(configPath string) (Config, error) {
	var file *os.File
	file, _ = os.Open(configPath)

	v := viper.New()
	v.SetConfigType("yaml")
	v.AutomaticEnv()

	/* default */
	v.SetDefault("app_id", "")
	v.SetDefault("release", false)
	v.SetDefault("logger.level", "INFO")
	v.SetDefault("logger.format", logger.ConsoleFormat)
	v.SetDefault("http.port", "8080")

	/* database */
	v.SetDefault("database.driver", "mysql")
	v.SetDefault("database.host", "localhost")
	v.SetDefault("database.port", 3306)
	v.SetDefault("database.database", "")
	v.SetDefault("database.instance_name", "")
	v.SetDefault("database.user", "root")
	v.SetDefault("database.password", "")
	v.SetDefault("database.connect_timeout", "10s")
	v.SetDefault("database.read_timeout", "30s")
	v.SetDefault("database.write_timeout", "30s")
	v.SetDefault("database.dial_timeout", "10s")
	v.SetDefault("database.max_idletime", "1h")
	v.SetDefault("database.max_lifetime", "1h")
	v.SetDefault("database.max_idle_conn", 2)
	v.SetDefault("database.max_open_conn", 5)
	v.SetDefault("database.ssl_mode", false)

	/* mq */
	v.SetDefault("mq.driver", "")
	/* watermill kafka option */
	v.SetDefault("mq.kafka_option.brokers", []string{})
	v.SetDefault("mq.kafka_option.consumer_group", "")
	v.SetDefault("mq.kafka_option.offsets_initial", -2) // OffsetNewest = -1, OffsetOldest = -2
	v.SetDefault("mq.kafka_option.fetch_default", 1024*1024)
	v.SetDefault("mq.kafka_option.required_acks", 1) // NoResponse = 0, WaitForLocal = 1, WaitForAll = -1
	/* confluent kafka option */
	v.SetDefault("mq.confluentkafka_option.brokers", []string{})
	v.SetDefault("mq.confluentkafka_option.client_id", "")
	v.SetDefault("mq.confluentkafka_option.consumer_group", "")
	v.SetDefault("mq.confluentkafka_option.group_id", "")
	v.SetDefault("mq.confluentkafka_option.offsets_initial", "")
	v.SetDefault("mq.confluentkafka_option.rebalance_strategy", "")
	v.SetDefault("mq.confluentkafka_option.group_instance_id", "")
	v.SetDefault("mq.confluentkafka_option.heartbeat_interval_ms", 0)
	v.SetDefault("mq.confluentkafka_option.session_timeout_ms", 0)
	v.SetDefault("mq.confluentkafka_option.auto_commit_interval_ms", 0)
	v.SetDefault("mq.confluentkafka_option.enable_auto_commit", false)
	v.SetDefault("mq.confluentkafka_option.acks", 2)
	v.SetDefault("mq.confluentkafka_option.compression_type", "gzip")
	v.SetDefault("mq.confluentkafka_option.retries", 5)
	v.SetDefault("mq.confluentkafka_option.batch_size", 0)
	v.SetDefault("mq.confluentkafka_option.flush_wait_ms", 0)
	v.SetDefault("mq.confluentkafka_option.fetch_max_bytes", 0)
	v.SetDefault("mq.confluentkafka_option.max_partition_fetch_bytes", 0)
	v.SetDefault("mq.confluentkafka_option.poll_timeout_ms", 100)
	v.SetDefault("mq.confluentkafka_option.sasl_username", "")
	v.SetDefault("mq.confluentkafka_option.sasl_password", "")
	v.SetDefault("mq.confluentkafka_option.sasl_mechanisms", "")
	v.SetDefault("mq.confluentkafka_option.security_protoco", "")

	/* eth client */
	v.SetDefault("eth_client.url", "")
	v.SetDefault("eth_client.dial_timeout", 10*time.Second)
	v.SetDefault("eth_client.max_client_conn", 100)

	/* scheduler */
	v.SetDefault("scheduler.unstable_num", 20)
	v.SetDefault("scheduler.start_at", 0)
	v.SetDefault("scheduler.sync.interval", 10*time.Second)
	v.SetDefault("scheduler.batch_limit", 100)

	/* crawler */
	v.SetDefault("crawler.topic", "")
	v.SetDefault("crawler.pool_size", 200)
	v.SetDefault("crawler.timeout", 10*time.Second)

	/* database writer */
	v.SetDefault("database_writer.topic", "")
	v.SetDefault("database_writer.pool_size", 200)
	v.SetDefault("database_writer.timeout", 10*time.Second)

	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.ReadConfig(file)

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return config, err
	}

	return config, nil
}
