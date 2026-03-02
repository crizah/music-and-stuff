package models

type Playlist struct {
	PlaylistId string
	UserId     string
	Songs      []Song

	// Likes      int
	// Dislikes   int
	// Display    []Display
	// Comments   []Comments
}
type Card struct {
	UserId          string
	Saved           []string // array of uids who saved playlist (can be removed)
	Likes           int
	Dislikes        int
	Neutral         int
	Description     string
	PlaylistDisplay DisplayItem
	Tags            []Tag
	Comments        []Comment
}

type DisplayItem struct {
	Albums      []string
	Artists     []string
	AlbumCovers []string // s3 key ids or key id where album covers are stored
}

type Song struct {
	SongID string
	Album  string
	Artist string
	Tag    Tag
}

type Tag struct {
	Genere               string
	Vibe                 string
	MoreAlgoSpecificInfo string
}
