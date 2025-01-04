package main

import (
	"backend-bootcamp-assignment-2024/internal/auth"
	"backend-bootcamp-assignment-2024/internal/config"
	"backend-bootcamp-assignment-2024/internal/db"
	"backend-bootcamp-assignment-2024/internal/routes"
	"fmt"
	"log"
	"net/http"
)

func main() {
	cfg := config.NewConfig()
	db, err := db.NewDB(cfg.DBConnStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	authService := auth.NewAuthService([]byte(cfg.SecretKey))
	authHandler := auth.NewAuthHandler(authService)

	router := routes.SetupRouter(authHandler)

	fmt.Println("Server is running on port 8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}
