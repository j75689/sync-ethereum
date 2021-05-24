# Introduction

`sync-ethereum` is project that can help you synchronize Ethereum's block record and transaction records and write them to the database.


## Usage
```
Usage:
   [command]

Available Commands:
  crawler     Start crawler
  help        Help about any command
  http        Start http server
  migrate     Migration tool
  scheduler   Start scheduler
  writer      Start database writer

Flags:
      --config string   config file (default "config/default.config.yaml")
  -h, --help            help for this command
      --timeout uint    graceful shutdown timeout (second) (default 300)

Use " [command] --help" for more information about a command.
```

## Configuration

| name | env | type | option | desc | default|
|---|---|---|---|---|---|
| app_id | APP_ID | string | | application name | `""`|
| release | RELEASE | bool | | is it a release version | `false` |
| logger.level | LOGGER_LEVEL | string | `ERROR`、`WARN`、`INFO`、`DEBUG`、`TRACE` | log level | `INFO` |
| logger.format | LOGGER_FORMAT | string | `console`、`json` | log format | `console` |
| http.port | HTTP_PORT | int | | http port | `8080` |
|---|---|---|---|---|---|
| database.driver | DATABASE_DRIVER | string | `mysql`、`postgres`、`sqlite` | sql driver | `mysql` |
| database.host | DATABASE_HOST | string | | database host | `""` |
| database.port | DATABASE_PORT | int | | database port | `3306` |
| database.user | DATABASE_USER | string | | database user | `""` |
| database.password | DATABASE_PASSWORD | int | | database password | `""` |
| database.database | DATABASE_DATABASE | int | | database name | `""` |
| database.max_open_conn | DATABASE_MAX_OPEN_CONN | int | | max open connection of database | `5` |
| database.min_idle_conn | DATABASE_MIN_IDLE_CONN | int | | max idle connection of database | `2` |
| database.connect_timeout | DATABASE_CONNECT_TIMEOUT | string | | connect database timeout | `10s` |
| database.read_timeout | DATABASE_READ_TIMEOUT | string | | read database timeout | `30s` |
| database.write_timeout | DATABASE_WRITE_TIMEOUT | string | | write database timeout | `30s` |
| database.dial_timeout | DATABASE_DIAL_TIMEOUT | time.duration | | ping database timeout | `10s` |
| database.max_idletime | DATABASE_MAX_IDLETIME | time.duration | | maximum amount of time a connection may be idle | `1h` |
| database.max_lifetime | DATABASE_MAX_LIFETIME | time.duration | | maximum amount of time a connection may be reused | `1h` |
| database.ssl_mode | DATABASE_SSL_MODE | bool | | connect database with ssl | `false` |
|---|---|---|---|---|---|
| mq.driver | MQ_DRIVER | string | `confluentkafka` | message queue driver | `""` |
| mq.confluentkafka_option.brokers | MQ_CONFLUENTKAFKA_OPTION_BROKERS | []string | | kafka broker list | `""` |
| mq.confluentkafka_option.consumer_group | MQ_CONFLUENTKAFKA_OPTION_CONSUMER_GROUP | string | | consumer group name | `""` |
| mq.confluentkafka_option.group_id | MQ_CONFLUENTKAFKA_OPTION_GROUP_ID | string | | consumer group id | `""` |
| mq.confluentkafka_option.client_id | MQ_CONFLUENTKAFKA_OPTION_CLIENT_ID | string | | client id | `""` |
| mq.confluentkafka_option.poll_timeout_ms | MQ_CONFLUENTKAFKA_POLL_TIMEOUT_MS | int | | millisecond of poll message | `100` |
|---|---|---|---|---|---|
| scheduler.unstable_num | SCHEDULER_UNSTABLE_NUM | string | | the latest quantity will be marked as unstable | `20` |
| scheduler.start_at | SCHEDULER_START_AT | int | | start synchronization from the block number | `0` |
| scheduler.batch_limit | SCHEDULER_BATCH_LIMIT | int | | limit of each synchronization | `100` |
| scheduler.sync.interval | SCHEDULER_SYNC_INTERVAL | time.duration | | interval of synchronization | `"10s"` |
|---|---|---|---|---|---|
| crawler.topic | CRAWLER_TOPIC | string | | topic name of the received message | `""` |
| crawler.pool_size | CRAWLER_POOL_SIZE | int | | worker size of crawler | `"200"` |
| crawler.timeout | CRAWLER_TIMEOUT | time.duration | | timeout of each operation | `10s` |
|---|---|---|---|---|---|
| writer.topic | WRITER_TOPIC | string | | topic name of the received message | `""` |
| writer.pool_size | WRITER_POOL_SIZE | int | | worker size of writer | `"200"` |
| writer.timeout | WRITER_TIMEOUT | time.duration | | timeout of each operation | `10s` |

## Docker Compose
[Example](https://github.com/j75689/sync-ethereum/blob/main/deployment/docker-compose/docker-compose.yaml)
```bash
make docker-compose-up
make docker-compose-down
```
