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
	Driver      string            `mapstructure:"driver"`
	KafkaOption KafkaOptionConfig `mapstructure:"kafka_option"`
}

type KafkaOptionConfig struct {
	Brokers        []string `mapstructure:"brokers"`
	ConsumerGroup  string   `mapstructure:"consumer_group"`
	OffsetsInitial int64    `mapstructure:"offsets_initial"`
}

type EthClientConfig struct {
	URL         string        `mapstructure:"url"`
	DialTimeout time.Duration `mapstructure:"dial_timeout"`
}

type SchedulerConfig struct {
	UnstableNumber int        `mapstructure:"unstable_num"`
	StartAt        int64      `mapstructure:"start_at"`
	Sync           SyncConfig `mapstructure:"sync"`
}
type SyncConfig struct {
	Interval time.Duration `mapstructure:"interval"`
}
type CrawlerConfig struct {
	Topic    string `mapstructure:"topic"`
	PoolSize int    `mapstructure:"pool_size"`
}

type DatabaseWriterConfig struct {
	Topic    string `mapstructure:"topic"`
	PoolSize int    `mapstructure:"pool_size"`
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
	v.SetDefault("mq.kafka_option.brokers", []string{})
	v.SetDefault("mq.kafka_option.consumer_group", "")
	v.SetDefault("mq.kafka_option.offsets_initial", -2) // OffsetNewest = -1 ,OffsetOldest = -2

	/* eth client */
	v.SetDefault("eth_client.url", "")
	v.SetDefault("eth_client.dial_timeout", 10*time.Second)

	/* scheduler */
	v.SetDefault("scheduler.unstable_num", 20)
	v.SetDefault("scheduler.start_at", 0)
	v.SetDefault("scheduler.sync.interval", 10*time.Second)

	/* crawler */
	v.SetDefault("crawler.topic", "")
	v.SetDefault("crawler.pool_size", 200)

	/* database writer */
	v.SetDefault("database_writer.topic", "")
	v.SetDefault("database_writer.pool_size", 200)

	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.ReadConfig(file)

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return config, err
	}

	return config, nil
}
