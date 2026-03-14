package main

import (
	"log"
	"server/internal/auth"
	"server/internal/middleware"
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
	r.Use(middleware.CORSMiddleware())

	v1 := r.Group("api/v1")
	v1.POST("/auth/signup", auth.SignUpHandler)
	v1.POST("/auth/login", auth.LoginHandler)

}
