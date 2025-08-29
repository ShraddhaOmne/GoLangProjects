package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"food-delivery/models"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/twmb/franz-go/pkg/kadm"
	"github.com/twmb/franz-go/pkg/kerr"
	"github.com/twmb/franz-go/pkg/kgo"
)

type Messaging struct {
	ChMessaging chan []byte
	Topic       string
	Brokers     []string
}

func NewMessaging(topic string, brokers []string) *Messaging {
	// Trim spaces from brokers to avoid connection issues
	for i := range brokers {
		brokers[i] = strings.TrimSpace(brokers[i])
	}

	// Create admin client
	cl, err := kgo.NewClient(kgo.SeedBrokers(brokers...))
	if err != nil {
		panic(fmt.Sprintf("failed to init kafka client: %v", err))
	}
	defer cl.Close()

	admin := kadm.NewClient(cl)

	ctx := context.Background()
	res, err := admin.CreateTopics(ctx, 3, 3, nil, topic) // 3 partitions, RF=1
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create topic")
	}
	for t, d := range res {
		if d.Err != nil {
			if d.Err == kerr.TopicAlreadyExists { // ✅ ignore if already exists
				log.Info().Str("topic", t).Msg("topic already exists, continuing...")
				continue
			}
			log.Fatal().Err(d.Err).Str("topic", t).Msg("topic creation failed")
		} else {
			log.Info().Str("topic", t).Msg("topic created successfully")
		}
	}
	log.Info().Str("topic", topic).Msg("topic ensured")

	return &Messaging{make(chan []byte), topic, brokers}
}

func (msg *Messaging) ProduceRecords() {
	fmt.Println("Topic : ", msg.Topic)
	if msg.Topic == "" {
		panic("invalid topic")
	}
	if len(msg.Brokers) < 1 {
		panic("invalid brokers")
	}

	cl, err := kgo.NewClient(
		kgo.SeedBrokers(msg.Brokers...),
		kgo.RequiredAcks(kgo.AllISRAcks()),
	)
	if err != nil {
		panic(err)
	}
	defer cl.Close()

	ctx := context.Background()
	for message := range msg.ChMessaging {
		record := &kgo.Record{Topic: msg.Topic, Value: message}

		cl.Produce(ctx, record, func(r *kgo.Record, err error) {
			if err != nil {
				fmt.Printf("record had a produce error: %v\n", err)
				return
			}
			order := new(models.Orders)
			_ = json.Unmarshal(r.Value, order)
			fmt.Println("Producer-->", r.ProducerID,
				"Topic-->", r.Topic,
				"Partition:", r.Partition,
				"Offset:", r.Offset,
				"Value:", order)
		})
	}
	cl.Flush(ctx)
	log.Print("Closed publishing data")
}

func (msg *Messaging) ConsumeRecords() {
	cl, err := kgo.NewClient(
		kgo.SeedBrokers(msg.Brokers...),
		// set a group id so multiple consumers can share load
		kgo.ConsumerGroup("food-delivery-service"),
		kgo.ConsumeTopics(msg.Topic),
	)
	if err != nil {
		panic(err)
	}
	defer cl.Close()

	ctx := context.Background()
	for {
		// Poll blocks until new messages or timeout
		fetches := cl.PollFetches(ctx)
		if errs := fetches.Errors(); len(errs) > 0 {
			for _, e := range errs {
				log.Error().Err(e.Err).Str("topic", e.Topic).Msg("consume error")
			}
			continue
		}

		// Iterate over all records in the fetch batch
		fetches.EachRecord(func(r *kgo.Record) {
			order := new(models.Orders)
			if err := json.Unmarshal(r.Value, order); err != nil {
				log.Error().Err(err).Str("topic", r.Topic).Msg("failed to unmarshal order")
				return
			}
			fmt.Printf("Consumer received: OrderID=%d, Status=%s, Partition=%d, Offset=%d\n",
				order.OrderId, order.Status, r.Partition, r.Offset)
		})

		// (optional) small sleep to avoid busy-looping
		time.Sleep(100 * time.Millisecond)
	}
}
