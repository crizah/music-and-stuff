package schemas

type SpotifyUser struct {
	ID           string `json:"id"`
	DisplayName  string `json:"display_name"`
	Email        string `json:"email"`
	ExternalURLs struct {
		Spotify string `json:"spotify"`
	} `json:"external_urls"`
}

type SpotifyTracksResponse struct {
	HREF     string `json:"href"`
	Limit    int    `json:"limit"`
	Next     string `json:"next"`
	Offset   int    `json:"offset"`
	Previous string `json:"previous"`
	Total    int    `json:"total"`
	Items    []struct {
		AddedAt string `json:"added_at"`
		Track   struct {
			ID    string `json:"id"`
			Name  string `json:"name"`
			Album struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			} `json:"album"`
			Artists []struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			} `json:"artists"`
			ExternalIDs struct {
				ISRC string `json:"isrc"`
			} `json:"external_ids"`
			Popularity int `json:"popularity"`
		} `json:"track"`
	} `json:"items"`
}

// Update your existing SpotifyPlaylistResponse if needed
type SpotifyPlaylistResponse struct {
	HREF     string `json:"href"`
	Limit    int    `json:"limit"`
	Next     string `json:"next"`
	Offset   int    `json:"offset"`
	Previous string `json:"previous"`
	Total    int    `json:"total"`
	Items    []struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		Public      bool   `json:"public"`
		Tracks      struct {
			Total int `json:"total"`
		} `json:"tracks"`
		ExternalURLs struct {
			Spotify string `json:"spotify"`
		} `json:"external_urls"`
		Images []struct {
			URL    string `json:"url"`
			Height int    `json:"height"`
			Width  int    `json:"width"`
		} `json:"images"`
	} `json:"items"`
}
