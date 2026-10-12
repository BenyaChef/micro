package queueservermodel

type ConsumerMessage struct {
	Topic     string
	Key       string
	Value     []byte
	Partition int32
	Offset    int64
}
