package repository

import (
	"context"

	"spotistats/internal/database"
)

type Stat struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

func GetGenreStats(ctx context.Context) ([]Stat, error) {
	query := `
		SELECT genre, COUNT(*) as count 
		FROM tracks 
		WHERE genre != '' 
		GROUP BY genre 
		ORDER BY count DESC 
		LIMIT 10
	`
	return queryStats(ctx, query)
}

func GetArtistStats(ctx context.Context) ([]Stat, error) {
	query := `
		SELECT artist, COUNT(*) as count 
		FROM tracks 
		GROUP BY artist 
		ORDER BY count DESC 
		LIMIT 10
	`
	return queryStats(ctx, query)
}

func GetYearStats(ctx context.Context) ([]Stat, error) {
	query := `
		SELECT year::text, COUNT(*) as count 
		FROM tracks 
		WHERE year > 0 
		GROUP BY year 
		ORDER BY year DESC 
		LIMIT 20
	`
	return queryStats(ctx, query)
}

func queryStats(ctx context.Context, query string) ([]Stat, error) {
	rows, err := database.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []Stat
	for rows.Next() {
		var s Stat
		if err := rows.Scan(&s.Name, &s.Count); err != nil {
			return nil, err
		}
		stats = append(stats, s)
	}
	return stats, nil
}
