package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"users-service/models"

	"github.com/rs/zerolog/log"
	"github.com/twmb/franz-go/pkg/kgo"
)

type Messaging struct {
	ChMessaging chan []byte
	Topic       string
	Brokers     []string
}

func NewMessaging(topic string, brokers []string) *Messaging {
	return &Messaging{make(chan []byte), topic, brokers}
}

func (msg *Messaging) ProduceRecords() {
	if msg.Topic == "" {
		panic("invalid topc")
	}

	if len(msg.Brokers) < 1 {
		panic("invalid brokers")
	}

	cl, err := kgo.NewClient(
		kgo.SeedBrokers(msg.Brokers...),
		// kgo.ConsumerGroup(ConsumerGroup),
		// kgo.ConsumeTopics(Topic),
		kgo.RequiredAcks(kgo.AllISRAcks()), // or kgo.RequireOneAck(), kgo.RequireNoAck()
		//kgo.DisableIdempotentWrite()
	)
	if err != nil {
		panic(err)
	}

	defer cl.Close()
	ctx := context.Background()
	for message := range msg.ChMessaging {
		record := &kgo.Record{Topic: msg.Topic, Value: message, Key: nil}
		cl.Produce(ctx, record, func(r *kgo.Record, err error) {
			//defer wg.Done()
			if err != nil {
				fmt.Printf("record had a produce error: %v\n", err)
			}
			user := new(models.User)
			json.Unmarshal(r.Value, user)
			//fmt.Println(user)
			fmt.Println("Producer-->", r.ProducerID, "Topid-->", r.Topic, "Partition:", r.Partition, "Offset:", r.Offset, "Value:", user)
		})

		// Alternatively, ProduceSync exists to synchronously produce a batch of records.
		if err := cl.ProduceSync(ctx, record).FirstErr(); err != nil {
			fmt.Printf("record had a produce error while synchronously producing: %v\n", err)
		}
	}
	cl.Flush(ctx)
	log.Print("Closed publishing data")

}
