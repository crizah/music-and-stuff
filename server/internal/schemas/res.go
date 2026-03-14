package schemas

type SignUpRes struct {
	Username     string `json:"username"`
	SessionToken string `json:"sessionToken"`
}
