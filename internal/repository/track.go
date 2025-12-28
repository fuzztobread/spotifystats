package repository

import (
	"context"

	"spotistats/internal/database"
	"spotistats/internal/models"
)

func InsertTrack(ctx context.Context, t models.Track) error {
	query := `
		INSERT INTO tracks (
			track_id, name, artist, spotify_preview_url, spotify_id,
			tags, genre, year, duration_ms, danceability, energy,
			key, loudness, mode, speechiness, acousticness,
			instrumentalness, liveness, valence, tempo, time_signature
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10,
			$11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21
		) ON CONFLICT (track_id) DO NOTHING
	`

	_, err := database.Pool.Exec(ctx, query,
		t.TrackID, t.Name, t.Artist, t.SpotifyPreviewURL, t.SpotifyID,
		t.Tags, t.Genre, t.Year, t.DurationMS, t.Danceability, t.Energy,
		t.Key, t.Loudness, t.Mode, t.Speechiness, t.Acousticness,
		t.Instrumentalness, t.Liveness, t.Valence, t.Tempo, t.TimeSignature,
	)
	return err
}

func InsertTracks(ctx context.Context, tracks []models.Track) error {
	for _, t := range tracks {
		if err := InsertTrack(ctx, t); err != nil {
			return err
		}
	}
	return nil
}

func GetTracks(ctx context.Context, limit, offset int, query, genre, artist string) ([]models.Track, int, error) {
	// count total
	countSQL := `SELECT COUNT(*) FROM tracks WHERE 1=1`
	args := []interface{}{}
	argIndex := 1

	if query != "" {
		countSQL += ` AND LOWER(name) LIKE LOWER($` + itoa(argIndex) + `)`
		args = append(args, "%"+query+"%")
		argIndex++
	}
	if genre != "" {
		countSQL += ` AND LOWER(genre) LIKE LOWER($` + itoa(argIndex) + `)`
		args = append(args, "%"+genre+"%")
		argIndex++
	}
	if artist != "" {
		countSQL += ` AND LOWER(artist) LIKE LOWER($` + itoa(argIndex) + `)`
		args = append(args, "%"+artist+"%")
		argIndex++
	}

	var total int
	err := database.Pool.QueryRow(ctx, countSQL, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// fetch tracks
	selectSQL := `
		SELECT track_id, name, artist, spotify_preview_url, spotify_id,
			tags, genre, year, duration_ms, danceability, energy,
			key, loudness, mode, speechiness, acousticness,
			instrumentalness, liveness, valence, tempo, time_signature
		FROM tracks WHERE 1=1
	`

	args = []interface{}{}
	argIndex = 1

	if query != "" {
		selectSQL += ` AND LOWER(name) LIKE LOWER($` + itoa(argIndex) + `)`
		args = append(args, "%"+query+"%")
		argIndex++
	}
	if genre != "" {
		selectSQL += ` AND LOWER(genre) LIKE LOWER($` + itoa(argIndex) + `)`
		args = append(args, "%"+genre+"%")
		argIndex++
	}
	if artist != "" {
		selectSQL += ` AND LOWER(artist) LIKE LOWER($` + itoa(argIndex) + `)`
		args = append(args, "%"+artist+"%")
		argIndex++
	}

	selectSQL += ` ORDER BY name LIMIT $` + itoa(argIndex) + ` OFFSET $` + itoa(argIndex+1)
	args = append(args, limit, offset)

	rows, err := database.Pool.Query(ctx, selectSQL, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var tracks []models.Track
	for rows.Next() {
		var t models.Track
		err := rows.Scan(
			&t.TrackID, &t.Name, &t.Artist, &t.SpotifyPreviewURL, &t.SpotifyID,
			&t.Tags, &t.Genre, &t.Year, &t.DurationMS, &t.Danceability, &t.Energy,
			&t.Key, &t.Loudness, &t.Mode, &t.Speechiness, &t.Acousticness,
			&t.Instrumentalness, &t.Liveness, &t.Valence, &t.Tempo, &t.TimeSignature,
		)
		if err != nil {
			return nil, 0, err
		}
		tracks = append(tracks, t)
	}

	return tracks, total, nil
}

func GetTrackByID(ctx context.Context, id string) (*models.Track, error) {
	query := `
		SELECT track_id, name, artist, spotify_preview_url, spotify_id,
			tags, genre, year, duration_ms, danceability, energy,
			key, loudness, mode, speechiness, acousticness,
			instrumentalness, liveness, valence, tempo, time_signature
		FROM tracks WHERE track_id = $1 OR spotify_id = $1
	`

	var t models.Track
	err := database.Pool.QueryRow(ctx, query, id).Scan(
		&t.TrackID, &t.Name, &t.Artist, &t.SpotifyPreviewURL, &t.SpotifyID,
		&t.Tags, &t.Genre, &t.Year, &t.DurationMS, &t.Danceability, &t.Energy,
		&t.Key, &t.Loudness, &t.Mode, &t.Speechiness, &t.Acousticness,
		&t.Instrumentalness, &t.Liveness, &t.Valence, &t.Tempo, &t.TimeSignature,
	)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func itoa(i int) string {
	return string(rune('0' + i))
}
