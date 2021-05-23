package confluentkafka

type KafkaOption struct {
	Brokers       []string
	ClientID      string
	ConsumerGroup string
	// All clients sharing the same group.id belong to the same group.
	GroupID string
	// earliest / latest
	OffsetsInitial string
	// range / roundrobin / cooperative-sticky; default cooperative-sticky
	RebalanceStrategy    string
	GroupInstanceID      string
	HeartbeatIntervalMs  int
	SessionTimeoutMs     int
	AutoCommitIntervalMs int
	EnableAutoCommit     bool

	// ================ producer related config ================
	// valid value: 1、2、3;
	// 1:
	// The producer will not wait for any acknowledgment from the server at all.
	// The record will be immediately added to the socket buffer and considered sent.
	// 2:
	// This will mean the leader will write the record to its local log but will respond without awaiting full acknowledgement from all followers.
	// 3:
	// This means the leader will wait for the full set of in-sync replicas (ISR) to acknowledge the record.
	Acks int
	// none, gzip, snappy, lz4, zstd
	CompressionType string
	Retries         int
	BatchSize       int
	FlushWaitMs     int

	// ================ consumer related config ================
	FetchMaxBytes          int
	MaxPartitionFetchBytes int
	PollTimeoutMs          int
}
