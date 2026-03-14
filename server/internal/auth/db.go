package auth

import (
	"context"
	"errors"
	"time"

	"server/internal/models"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

func (s *Server) addToUsers(username string, user_id string) error {
	now := time.Now()
	user := models.Users{
		Id:        user_id,
		Username:  username,
		CreatedAt: now.GoString(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := s.CollUsers.InsertOne(ctx, user)
	if err != nil {
		// Duplicate key error (username)
		if mongo.IsDuplicateKeyError(err) {
			return errors.New("user already exists")
		}
		return err
	}

	return nil

}
func (s *Server) addToPasswords(user_id string, salt string, hash string) error {
	pass := models.Passwords{
		Id:     user_id,
		Hashed: hash,
		Salt:   salt,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := s.CollPasswords.InsertOne(ctx, pass)
	if err != nil {
		// Duplicate key error
		if mongo.IsDuplicateKeyError(err) {
			return errors.New("user already exists")
		}
		return err
	}

	return nil

}

func (s *Server) addToSpotify(user_id string, id string, playlists []models.Playlist) error {
	spotUser := models.SpotifyClient{
		Id:        user_id,
		SpotifyId: id,
		Playlists: playlists,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := s.CollSpotify.InsertOne(ctx, spotUser)
	if err != nil {
		// Duplicate key error
		if mongo.IsDuplicateKeyError(err) {
			return errors.New("user already exists")
		}
		return err
	}

	return nil

}
