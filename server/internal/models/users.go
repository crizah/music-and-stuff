package models

type Users struct {
	Id        string `bson:"_id" json:"id"`
	Username  string `bson:"username" json:"username"`
	CreatedAt string `bson:"createdat" json:"createdat"`
}

type Passwords struct {
	Id     string `bson:"_id" json:"id"` // user.id
	Hashed string `bson:"hashed" json:"hash"`
	Salt   string `bson:"salt" json:"salt"`
}

type SpotifyClient struct {
	Id        string     `bson:"_id" json:"id"` // user.id
	SpotifyId string     `bson:"spotifyid" json:"spotifyid"`
	Playlists []Playlist `bson:"playlists" json:"playlists"`
}
