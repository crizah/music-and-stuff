package models

type Playlist struct {
	PlaylistId string `bson:"_id" json:"playlistid"`
	Songs      []Song `bson:"songs" json:"songs"`
}

type Song struct {
	SongID  string   `bson:"_id" json:"songid"`
	Album   string   `bson:"alubum" json:"album"`
	Artists []string `bson:"artists" json:"artists"`
	Tags    []Tag    `bson:"tags" json:"tags"`
}

type Tag struct {
	Genere               string
	Vibe                 string
	MoreAlgoSpecificInfo string
}
