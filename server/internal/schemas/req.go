package schemas

type SignUpReq struct {
	Username  string `json:"username" binding:"required,min=1,max=10"`
	Password  string `json:"password" binding:"required,min=1,max=20"`
	SpotifyId string `json:"spotifyid"`
}

type LoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
