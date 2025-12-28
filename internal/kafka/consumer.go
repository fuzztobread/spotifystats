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
		Brokers:  brokers,
		Topic:    topic,
		GroupID:  groupID,
		MinBytes: 1,
		MaxBytes: 10e6,
		MaxWait:  1 * time.Second,
	})

	log.Printf("kafka consumer started for topic: %s", topic)

	go func() {
		log.Println("consumer goroutine started, waiting for messages...")
		for {
			log.Println("waiting for next message...")
			msg, err := reader.ReadMessage(context.Background())
			if err != nil {
				log.Printf("[KAFKA READ ERROR] %v", err)
				time.Sleep(1 * time.Second)
				continue
			}

			log.Printf("received message: %s", string(msg.Value))

			var track models.Track
			if err := json.Unmarshal(msg.Value, &track); err != nil {
				log.Printf("[KAFKA UNMARSHAL ERROR] %v", err)
				continue
			}

			if err := repository.InsertTrack(context.Background(), track); err != nil {
				log.Printf("[DB INSERT ERROR] %v", err)
				continue
			}

			log.Printf("[KAFKA CONSUMED] track=%s artist=%s", track.Name, track.Artist)
		}
	}()
}
