package main

import (
	"log"
	"server/internal/auth"
	"server/internal/server"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	base, err := server.InitialiseServer()

	if err != nil {
		log.Fatal(err)
	}

	auth := auth.NewAuthServer(base)
	r := gin.Default()
	// r.Use(CORSMiddleware())

	v1 := r.Group("/v1")
	v1.POST("/auth/register", auth.SignUp) // works

}
