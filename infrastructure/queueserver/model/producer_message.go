package queueservermodel

type ProducerMessage struct {
	Topic        string
	PartitionKey string
	Value        []byte
}
