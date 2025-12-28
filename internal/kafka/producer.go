package kafka

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

var Writer *kafka.Writer

func InitProducer(brokers []string, topic string) {
	Writer = &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		BatchTimeout: 10 * time.Millisecond,
	}
	log.Printf("kafka producer initialized for topic: %s", topic)
}

func CloseProducer() {
	if Writer != nil {
		Writer.Close()
	}
}

func Publish(ctx context.Context, key string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	msg := kafka.Message{
		Key:   []byte(key),
		Value: data,
	}

	err = Writer.WriteMessages(ctx, msg)
	if err != nil {
		log.Printf("[KAFKA ERROR] %v", err)
		return err
	}

	log.Printf("[KAFKA PUBLISHED] key=%s", key)
	return nil
}
