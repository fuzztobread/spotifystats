package kafka

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"spotistats/internal/models"
	"spotistats/internal/repository"

	"github.com/segmentio/kafka-go"
)

func StartConsumer(brokers []string, topic, groupID string) {
	log.Printf("connecting consumer to brokers: %v", brokers)

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     brokers,
		Topic:       topic,
		GroupID:     groupID,
		StartOffset: kafka.LastOffset,
	})

	log.Printf("kafka consumer started for topic: %s", topic)

	go func() {
		log.Println("consumer goroutine started, waiting for messages...")
		for {
			msg, err := reader.FetchMessage(context.Background())
			if err != nil {
				log.Printf("[KAFKA FETCH ERROR] %v", err)
				time.Sleep(1 * time.Second)
				continue
			}

			log.Printf("received message key=%s", string(msg.Key))

			var track models.Track
			if err := json.Unmarshal(msg.Value, &track); err != nil {
				log.Printf("[KAFKA UNMARSHAL ERROR] %v", err)
				continue
			}

			if err := repository.InsertTrack(context.Background(), track); err != nil {
				log.Printf("[DB INSERT ERROR] %v", err)
				continue
			}

			if err := reader.CommitMessages(context.Background(), msg); err != nil {
				log.Printf("[KAFKA COMMIT ERROR] %v", err)
			}

			log.Printf("[KAFKA CONSUMED] track=%s artist=%s", track.Name, track.Artist)
		}
	}()
}
