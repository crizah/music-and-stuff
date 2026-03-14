package models

// type Card struct {
// 	UserId          string
// 	Saved           []string // array of uids who saved playlist (can be removed)
// 	Likes           int
// 	Dislikes        int
// 	Neutral         int
// 	Description     string
// 	PlaylistDisplay DisplayItem
// 	Tags            []Tag
// 	Comments        []Comment
// }

// type DisplayItem struct {
// 	Albums      []string
// 	Artists     []string
// 	AlbumCovers []string // s3 key ids or key id where album covers are stored
// }

type Playlist struct {
	PlaylistId string
	Songs      []Song
}

type Song struct {
	SongID  string
	Album   string
	Artists []string
	Tag     Tag
}

type Tag struct {
	Genere               string
	Vibe                 string
	MoreAlgoSpecificInfo string
}
