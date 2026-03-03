package server

import (
	"context"
	"os"
	"time"

	"github.com/go-redis/cache/v9"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// mongo and stuff
// redis

type Server struct {
	MongoClient   *mongo.Client
	JwtSecret     []byte
	CollUsers     *mongo.Collection
	CollPasswords *mongo.Collection
	CollSpotify   *mongo.Collection
	// CollAnswerLog *mongo.Collection
	StateCache *cache.Cache
}

func InitialiseServer() (*Server, error) {

	uri := os.Getenv("MONGO_DB_URI")
	token := []byte(os.Getenv("JWT_SECRET"))

	clientOptions := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(clientOptions)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	u := client.Database("Music").Collection("Users")
	p := client.Database("Music").Collection("Passwords")
	s := client.Database("Music").Collection("Spotify")

	// p.Indexes().CreateOne(ctx, mongo.IndexModel{
	// 	Keys: bson.D{{Key: "totalScore", Value: -1}},
	// })
	// p.Indexes().CreateOne(ctx, mongo.IndexModel{
	// 	Keys: bson.D{{Key: "maxStreak", Value: -1}},
	// })

	// a.Indexes().CreateOne(ctx, mongo.IndexModel{
	// 	Keys:    bson.D{{Key: "ikey", Value: 1}},
	// 	Options: options.Index().SetUnique(true),
	// })

	ring := redis.NewRing(&redis.RingOptions{
		Addrs: map[string]string{
			"server1": ":6379",
			"server2": ":6380",
		},
	})

	mycache := cache.New(&cache.Options{
		Redis:      ring,
		LocalCache: cache.NewTinyLFU(1000, time.Minute),
	})

	return &Server{MongoClient: client, CollUsers: u, CollPasswords: p,
		StateCache: mycache, JwtSecret: token, CollSpotify: s}, nil

}

func (s *Server) GenerateJWT(username string) (string, error) {
	claims := jwt.MapClaims{
		"sub": username,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(30 * 24 * time.Hour).Unix(), // 30 days
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.JwtSecret)

}
