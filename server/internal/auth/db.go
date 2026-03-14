package auth

import (
	"context"
	"errors"
	"time"

	"server/internal/models"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var (
	ErrUserNotFound = errors.New("user not found")
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

// func (s *Server) deleteUser(user_id string) error {
// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()

// 	_, err := s.CollUsers.DeleteOne(ctx, bson.M{"_id": user_id})
// 	return err
// }

// func (s *Server) deletePasswords(user_id string) error {
// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()

// 	_, err := s.CollPasswords.DeleteOne(ctx, bson.M{"_id": user_id})
// 	return err
// }

// func (s *Server) deleteSpotify(user_id string) error {
// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()

// 	_, err := s.CollSpotify.DeleteOne(ctx, bson.M{"_id": user_id})
// 	return err
// }

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

func (s *Server) getIdFromUsername(username string) (*models.Users, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user models.Users
	err := s.CollUsers.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (s *Server) getSaltAndHash(user_id string) (*models.Passwords, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var pass models.Passwords
	err := s.CollPasswords.FindOne(ctx, bson.M{"_id": user_id}).Decode(&pass)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &pass, nil

}
