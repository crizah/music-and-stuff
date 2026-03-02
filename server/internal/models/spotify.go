package models

type SpotifyClient struct {
	Id        string     `json:"id"`
	SpotifyId string     `json:"spotifyid"`
	Followers int        `json:"followers"`
	Playlists []Playlist `json:"playlists"`
}
