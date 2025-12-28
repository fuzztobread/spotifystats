package models

type Track struct {
	TrackID           string  `json:"track_id" csv:"track_id"`
	Name              string  `json:"name" csv:"name"`
	Artist            string  `json:"artist" csv:"artist"`
	SpotifyPreviewURL string  `json:"spotify_preview_url" csv:"spotify_preview_url"`
	SpotifyID         string  `json:"spotify_id" csv:"spotify_id"`
	Tags              string  `json:"tags" csv:"tags"`
	Genre             string  `json:"genre" csv:"genre"`
	Year              int     `json:"year" csv:"year"`
	DurationMS        int     `json:"duration_ms" csv:"duration_ms"`
	Danceability      float64 `json:"danceability" csv:"danceability"`
	Energy            float64 `json:"energy" csv:"energy"`
	Key               int     `json:"key" csv:"key"`
	Loudness          float64 `json:"loudness" csv:"loudness"`
	Mode              int     `json:"mode" csv:"mode"`
	Speechiness       float64 `json:"speechiness" csv:"speechiness"`
	Acousticness      float64 `json:"acousticness" csv:"acousticness"`
	Instrumentalness  float64 `json:"instrumentalness" csv:"instrumentalness"`
	Liveness          float64 `json:"liveness" csv:"liveness"`
	Valence           float64 `json:"valence" csv:"valence"`
	Tempo             float64 `json:"tempo" csv:"tempo"`
	TimeSignature     int     `json:"time_signature" csv:"time_signature"`
}
