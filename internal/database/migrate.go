package database

import (
	"context"
	"log"
)

func Migrate() error {
	query := `
	CREATE TABLE IF NOT EXISTS tracks (
		track_id VARCHAR(50) PRIMARY KEY,
		name TEXT NOT NULL,
		artist TEXT NOT NULL,
		spotify_preview_url TEXT,
		spotify_id VARCHAR(50),
		tags TEXT,
		genre VARCHAR(100),
		year INT,
		duration_ms INT,
		danceability FLOAT,
		energy FLOAT,
		key INT,
		loudness FLOAT,
		mode INT,
		speechiness FLOAT,
		acousticness FLOAT,
		instrumentalness FLOAT,
		liveness FLOAT,
		valence FLOAT,
		tempo FLOAT,
		time_signature INT
	);

	CREATE INDEX IF NOT EXISTS idx_tracks_artist ON tracks(artist);
	CREATE INDEX IF NOT EXISTS idx_tracks_genre ON tracks(genre);
	CREATE INDEX IF NOT EXISTS idx_tracks_year ON tracks(year);
	`
	_, err := Pool.Exec(context.Background(), query)
	if err != nil {
		return err
	}

	log.Println("migration completed")
	return nil
}
