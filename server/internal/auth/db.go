package auth

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func (s *Server) addToUsers(username string, user_id string) error {
	user := bson.M{
		"_id":       username,
		"username":  username,
		"createdAt": time.Now().UTC(),
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
	user := bson.M{
		"_id":    user_id, // PK
		"salt":   salt,
		"hashed": hash,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := s.CollPasswords.InsertOne(ctx, user)
	if err != nil {
		// Duplicate key error
		if mongo.IsDuplicateKeyError(err) {
			return errors.New("user already exists")
		}
		return err
	}

	return nil

}

func (s *Server) addToSpotify(user_id string, id string) error {
	user := bson.M{
		"_id":       user_id, // PK
		"spotifyid": id,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := s.CollSpotify.InsertOne(ctx, user)
	if err != nil {
		// Duplicate key error
		if mongo.IsDuplicateKeyError(err) {
			return errors.New("user already exists")
		}
		return err
	}

	return nil

}
