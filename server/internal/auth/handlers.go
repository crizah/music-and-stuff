package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type SignUpReq struct {
	Username string `json:"username" binding:"required,min=1,max=10"`
	Password string `json:"password" binding:"required,min=1,max=20"`
}

type SignUpRes struct {
	Username     string `json:"username"`
	SessionToken string `json:"sessionToken"`
}

func (s *Server) SignUp(c *gin.Context) {
	var req SignUpReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username bw 1-10 chars"})
		return

	}
	req.Username = strings.TrimSpace(req.Username)
	req.Password = strings.TrimSpace(req.Password)

	// put into db

}

// func (s *Server) RegisterUser(c *gin.Context) {
// 	var req RegisterReq
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "username bw 1-10 chars"})
// 		return
// 	}

// 	req.Username = strings.TrimSpace(req.Username)

// 	// put into db

// 	err := s.PutUserIntoDb(req.Username)
// 	if err != nil {
// 		if err.Error() == USER_EXISTS {
// 			c.JSON(http.StatusConflict, gin.H{"error": "username already taken"})
// 			return
// 		}

// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
// 		return

// 	}

// }
