package kafka

import (
	"context"
	"encoding/json"
	"log"

	"spotistats/internal/models"
	"spotistats/internal/repository"

	"github.com/segmentio/kafka-go"
)

func StartConsumer(brokers []string, topic, groupID string) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		Topic:    topic,
		GroupID:  groupID,
		MinBytes: 10e3,
		MaxBytes: 10e6,
	})

	log.Printf("kafka consumer started for topic: %s", topic)

	go func() {
		for {
			msg, err := reader.ReadMessage(context.Background())
			if err != nil {
				log.Printf("[KAFKA READ ERROR] %v", err)
				continue
			}

			var track models.Track
			if err := json.Unmarshal(msg.Value, &track); err != nil {
				log.Printf("[KAFKA UNMARSHAL ERROR] %v", err)
				continue
			}

			// insert to database
			if err := repository.InsertTrack(context.Background(), track); err != nil {
				log.Printf("[DB INSERT ERROR] %v", err)
				continue
			}

			log.Printf("[KAFKA CONSUMED] track=%s artist=%s", track.Name, track.Artist)
		}
	}()
}
