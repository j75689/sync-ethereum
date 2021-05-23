package confluentkafka

type MsgData struct {
	RequestID string       `json:"request_id"`
	Data      []byte       `json:"data,omitempty"`
	ConsumeID string       `json:"consume_id"`
	Commit    func() error `json:"-"`
}
