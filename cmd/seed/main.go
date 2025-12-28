package main

import (
	"context"
	"log"
	"time"

	"spotistats/internal/database"
	"spotistats/internal/loader"
	"spotistats/internal/repository"
)

func main() {
	// connect to postgres
	err := database.Connect("postgres://spotistats:spotistats@localhost:5432/spotistats")
	if err != nil {
		log.Fatal("failed to connect to db:", err)
	}
	defer database.Close()

	// run migrations
	if err := database.Migrate(); err != nil {
		log.Fatal("migration failed:", err)
	}

	// load CSV
	tracks, err := loader.LoadTracksFromCSV("data/tracks.csv")
	if err != nil {
		log.Fatal("failed to load CSV:", err)
	}
	log.Printf("loaded %d tracks from CSV", len(tracks))

	// seed to database
	log.Println("seeding to database...")
	start := time.Now()

	ctx := context.Background()
	if err := repository.InsertTracks(ctx, tracks); err != nil {
		log.Fatal("seeding failed:", err)
	}

	log.Printf("seeding completed in %v", time.Since(start))
}
