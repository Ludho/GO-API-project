package main

import (
	"fmt"
	"os"
	"trades/config"
	"trades/handlers"
	"trades/repos"
	"trades/services"

	"github.com/labstack/echo/v4"
)

func main() {
	server := echo.New()

	// load config
	config := config.Load()

	hashSecret := os.Getenv("SECRET_HASH")
	fmt.Printf("Secret_hash = %s", hashSecret)
	jwtSecret := os.Getenv("SECRET_JWT")

	userRepo := repos.NewUserRepository(config.Connection)
	userService := services.NewUserService(userRepo, hashSecret, jwtSecret)
	userHandler := handlers.NewUserHandler(userService)

	server.POST("/users", userHandler.Register)
	server.POST("/login", userHandler.Login)

	if err := server.Start(":8080"); err != nil {
		fmt.Println(err)
	}

}
