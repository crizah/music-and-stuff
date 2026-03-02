package models

type Users struct {
	Id        string `json:"id"`
	Username  string `json:"username"`
	CreatedAt string `json:"createdat"`
	Verified  bool   `json:"verified"`
}
