package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"server/internal/models"
	"server/internal/schemas"
	"time"
)

func (s *Server) VerifySpotifyUser(spotifyID string) (bool, error) {
	url := "https://api.spotify.com/v1/me"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", s.AccessToken))

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("spotify verification failed: %d", resp.StatusCode)
	}

	var user schemas.SpotifyUser
	err = json.NewDecoder(resp.Body).Decode(&user)
	if err != nil {
		return false, err
	}

	// Verify the ID matches
	if spotifyID != "" && user.ID != spotifyID {
		return false, nil
	}

	return true, nil
}

func (s *Server) GetUserPlaylists(spotifyID string) (*schemas.SpotifyPlaylistResponse, error) {
	url := fmt.Sprintf("https://api.spotify.com/v1/users/%s/playlists", spotifyID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", s.AccessToken))

	// Add query parameters for pagination
	q := req.URL.Query()
	q.Add("limit", "50")
	q.Add("offset", "0")
	req.URL.RawQuery = q.Encode()

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("spotify API error: %d - %s", resp.StatusCode, string(body))
	}

	var playlists schemas.SpotifyPlaylistResponse
	err = json.NewDecoder(resp.Body).Decode(&playlists)
	if err != nil {
		return nil, err
	}

	return &playlists, nil

}

func (s *Server) ExtractInfoFromPlaylists(spotifyRes *schemas.SpotifyPlaylistResponse) (*[]models.Playlist, error) {
	if spotifyRes == nil {
		return nil, fmt.Errorf("spotify response is nil")
	}

	if len(spotifyRes.Items) == 0 {
		return nil, fmt.Errorf("no playlists found")
	}

	playlists := make([]models.Playlist, 0)

	// per playlist
	for _, spotifyPlaylist := range spotifyRes.Items {
		playlist := models.Playlist{
			PlaylistId: spotifyPlaylist.ID,
			Songs:      []models.Song{},
		}

		// playlist tracks
		tracks, err := s.GetPlaylistTracks(spotifyPlaylist.ID)
		if err != nil {
			// Log error but continue with next playlist
			continue
		}

		// songs from tracks
		for _, track := range tracks.Items {
			song := models.Song{
				SongID:  track.Track.ID,
				Album:   track.Track.Album.Name,
				Artists: getArtistNames(track.Track.Artists),
			}

			// Generate tags for the song
			tag, err := s.GenerateTags(track.Track)
			if err != nil {
				// Log error but continue processing other songs
				continue
			}

			song.Tags = *tag
			playlist.Songs = append(playlist.Songs, song)
		}

		// Only add playlist if it has songs
		if len(playlist.Songs) > 0 {
			playlists = append(playlists, playlist)
		}
	}

	if len(playlists) == 0 {
		return nil, fmt.Errorf("no playlists with songs found")
	}

	return &playlists, nil
}

// all tracks from a Spotify playlist
func (s *Server) GetPlaylistTracks(playlistID string) (*schemas.SpotifyTracksResponse, error) {
	accessToken := s.AccessToken

	url := fmt.Sprintf("https://api.spotify.com/v1/playlists/%s/tracks", playlistID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	// Add query parameters
	q := req.URL.Query()
	q.Add("limit", "50")
	q.Add("offset", "0")
	req.URL.RawQuery = q.Encode()

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("spotify API error: %d - %s", resp.StatusCode, string(body))
	}

	var tracks schemas.SpotifyTracksResponse
	err = json.NewDecoder(resp.Body).Decode(&tracks)
	if err != nil {
		return nil, err
	}

	return &tracks, nil
}

// extract artist names from a list of artists
func getArtistNames(artists []struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}) []string {
	artistNames := make([]string, len(artists))
	for i, artist := range artists {
		artistNames[i] = artist.Name
	}
	return artistNames
}

func (s *Server) GenerateTags(track interface{}) (*[]models.Tag, error) {
	// TODO: Implement tag generation logic
	return nil, nil
}
