package loader

import (
	"os"

	"spotistats/internal/models"

	"github.com/gocarina/gocsv"
)

func LoadTracksFromCSV(filepath string) ([]models.Track, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var tracks []models.Track
	if err := gocsv.UnmarshalFile(file, &tracks); err != nil {
		return nil, err
	}

	return tracks, nil
}
