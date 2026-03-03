package auth

import (
	"net/http"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SignUpReq struct {
	Username  string `json:"username" binding:"required,min=1,max=10"`
	Password  string `json:"password" binding:"required,min=1,max=20"`
	SpotifyId string `json:"spotifyid"`
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

	// generate salt and hash
	salt, hash, err := HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "mysterious error: " + err.Error()})
		return
	}

	// create user id
	user_id := uuid.New().String()

	// get spotify info to add to potify table
	// verify spotify user accoiunt and add to users table

	// put username and id in users
	// put salt and hash in passwords
	// put spotify in spotify

	var wg sync.WaitGroup
	var mux sync.Mutex
	wg.Add(3)
	var errors struct {
		Errors []error
	}

	go func() {
		defer wg.Done()
		err := s.addToUsers(req.Username, user_id)
		mux.Lock()
		defer mux.Unlock()

		if err != nil {
			errors.Errors = append(errors.Errors, err)
		}

	}()

	go func() {
		defer wg.Done()
		err := s.addToPasswords(user_id, salt, hash)
		mux.Lock()
		defer mux.Unlock()

		if err != nil {
			errors.Errors = append(errors.Errors, err)
		}

	}()

	go func() {
		defer wg.Done()
		err := s.addToSpotify(user_id, req.SpotifyId)
		mux.Lock()
		defer mux.Unlock()

		if err != nil {
			errors.Errors = append(errors.Errors, err)
		}

	}()

	// create token

	token, err := s.GenerateJWT(req.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "token error"})
		return
	}

	c.JSON(http.StatusCreated, SignUpRes{
		Username:     req.Username,
		SessionToken: token,
	})

}
