package auth

import (
	"net/http"
	"strings"
	"sync"

	"server/internal/schemas"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (s *Server) SignUp(c *gin.Context) {
	var req schemas.SignUpReq
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

	// verify spotify user account
	verified, err := s.VerifySpotifyUser(req.SpotifyId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error" + err.Error()})

	}
	if !verified {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid spotify account"})
		return
	}

	// get users playlists

	playlistsRes, err := s.GetUserPlaylists(req.SpotifyId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to fetch playlists" + err.Error()})
		return
	}

	// extract relevent info from response
	playlists, err := s.ExtractInfoFromPlaylists(playlistsRes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "coudlnt extract playlist" + err.Error()})
	}

	// add to playlist

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
		err := s.addToSpotify(user_id, req.SpotifyId, *playlists)
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

	c.JSON(http.StatusCreated, schemas.SignUpRes{
		Username:     req.Username,
		SessionToken: token,
	})

}
