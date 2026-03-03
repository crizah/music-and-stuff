package models

type Users struct {
	Id        string `bson:"_id" json:"id"`
	Username  string `bson:"username" json:"username"`
	CreatedAt string `bson:"createdat" json:"createdat"`
	Verified  bool   `bson:"verified" json:"verified"`
}

type Passwords struct {
	Id     string `bson:"_id" json:"id"`
	Hashed string `bson:"hashed" json:"hash"`
	Salt   string `bson:"salt" json:"salt"`
}
