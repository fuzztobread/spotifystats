package models

type Track struct {
	TrackID           string  `json:"track_id"`
	Name              string  `json:"name"`
	Artist            string  `json:"artist"`
	SpotifyPreviewURL string  `json:"spotify_preview_url"`
	SpotifyID         string  `json:"spotify_id"`
	Tags              string  `json:"tags"`
	Genre             string  `json:"genre"`
	Year              int     `json:"year"`
	DurationMS        int     `json:"duration_ms"`
	Danceability      float64 `json:"danceability"`
	Energy            float64 `json:"energy"`
	Key               int     `json:"key"`
	Loudness          float64 `json:"loudness"`
	Mode              int     `json:"mode"`
	Speechiness       float64 `json:"speechiness"`
	Acousticness      float64 `json:"acousticness"`
	Instrumentalness  float64 `json:"instrumentalness"`
	Liveness          float64 `json:"liveness"`
	Valence           float64 `json:"valence"`
	Tempo             float64 `json:"tempo"`
	TimeSignature     int     `json:"time_signature"`
}
