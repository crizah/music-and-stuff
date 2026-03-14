package auth

import (
	"net/http"
	"strings"

	"server/internal/schemas"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (s *Server) SignUpHandler(c *gin.Context) {
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

	// var wg sync.WaitGroup
	// var mux sync.Mutex
	// wg.Add(3)
	// var errors struct {
	// 	Errors []error
	// }

	// go func() {
	// 	defer wg.Done()
	// 	err := s.addToUsers(req.Username, user_id)
	// 	mux.Lock()
	// 	defer mux.Unlock()

	// 	if err != nil {
	// 		errors.Errors = append(errors.Errors, err)
	// 	}

	// }()

	// go func() {
	// 	defer wg.Done()
	// 	err := s.addToPasswords(user_id, salt, hash)
	// 	mux.Lock()
	// 	defer mux.Unlock()

	// 	if err != nil {
	// 		errors.Errors = append(errors.Errors, err)
	// 	}

	// }()

	// go func() {
	// 	defer wg.Done()
	// 	err := s.addToSpotify(user_id, req.SpotifyId, *playlists)
	// 	mux.Lock()
	// 	defer mux.Unlock()

	// 	if err != nil {
	// 		errors.Errors = append(errors.Errors, err)
	// 	}

	// }()

	// dont do signup concurrently, partial creation might happen

	// stop at first error
	err = s.addToUsers(req.Username, user_id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = s.addToPasswords(user_id, salt, hash)
	if err != nil {
		// rollback

		// s.deleteUser(user_id)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = s.addToSpotify(user_id, req.SpotifyId, *playlists)
	if err != nil {
		// Rrollback
		// s.deleteUser(user_id)
		// s.deletePasswords(user_id)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

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

// add email and stuff for forgot passwoed and stuff later

func (s *Server) LoginHandler(c *gin.Context) {
	var req schemas.LoginReq

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return

	}
	req.Username = strings.TrimSpace(req.Username)
	req.Password = strings.TrimSpace(req.Password)

	// verify user exists
	user, err := s.getIdFromUsername(req.Username)
	if err != nil {
		if err == ErrUserNotFound {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "username or password wrong"})
			return

		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "error getting user" + err.Error()})
		return
	}

	// get user salt, hash

	password, err := s.getSaltAndHash(user.Id)
	if err != nil {
		if err == ErrUserNotFound {
			c.JSON(http.StatusConflict, gin.H{"error": "password for user not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "error getting user" + err.Error()})
		return

	}

	// verify password
	yay, err := VerifyPass(req.Password, password.Salt, password.Hashed)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "errorverifying password" + err.Error()})
		return

	}

	if !yay {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "username or password wrong" + err.Error()})
		return

	}

	// send success response

}
