package commands

import (
	"context"
	"log"
	"time"

	"spotistats/internal/config"
	"spotistats/internal/database"
	"spotistats/internal/loader"
	"spotistats/internal/repository"

	"github.com/spf13/cobra"
)

var seedFile string

var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "Seed database from CSV (skips if data exists)",
	Run:   runSeed,
}

func init() {
	seedCmd.Flags().StringVarP(&seedFile, "file", "f", "data/tracks.csv", "CSV file path")
}

func runSeed(cmd *cobra.Command, args []string) {
	cfg := config.Load()

	if err := database.Connect(cfg.DatabaseURL()); err != nil {
		log.Fatal("failed to connect to db:", err)
	}
	defer database.Close()

	// check if data already exists
	count, err := repository.CountTracks(context.Background())
	if err != nil {
		log.Fatal("failed to count tracks:", err)
	}

	if count > 0 {
		log.Printf("database already has %d tracks, skipping seed", count)
		return
	}

	// load CSV
	tracks, err := loader.LoadTracksFromCSV(seedFile)
	if err != nil {
		log.Fatal("failed to load CSV:", err)
	}
	log.Printf("loaded %d tracks from %s", len(tracks), seedFile)

	// seed to database
	log.Println("seeding to database...")
	start := time.Now()

	ctx := context.Background()
	if err := repository.InsertTracks(ctx, tracks); err != nil {
		log.Fatal("seeding failed:", err)
	}

	log.Printf("seeding completed in %v", time.Since(start))
}
