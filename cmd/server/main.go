package main

import (
	"backend-bootcamp-assignment-2024/internal/auth"
	"backend-bootcamp-assignment-2024/internal/config"
	"backend-bootcamp-assignment-2024/internal/db"
	"backend-bootcamp-assignment-2024/internal/flat"
	"backend-bootcamp-assignment-2024/internal/house"
	"backend-bootcamp-assignment-2024/internal/notification"
	"backend-bootcamp-assignment-2024/internal/routes"
	"backend-bootcamp-assignment-2024/pkg/sender"
	"context"
	"fmt"
	"log"
	"net/http"
)

func main() {
	cfg := config.NewConfig()
	db, err := db.NewDB(cfg.DBConnStr)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	ch, err := notification.InitRabbitMQ()
	if err != nil {
		log.Fatalf("failed to initialize RabbitMQ: %v", err)
	}
	defer ch.Close()

	publisher, err := notification.NewPublisher(ch)
	if err != nil {
		log.Fatalf("failed to create RabbitMQ publisher: %v", err)
	}

	authRepository := auth.NewAuthRepository(db)
	authService := auth.NewAuthService([]byte(cfg.SecretKey), authRepository)
	authHandler := auth.NewAuthHandler(authService)

	houseRepository := house.NewHouseRepository(db)
	houseService := house.NewHouseService(houseRepository)
	houseHandler := house.NewHouseHandler(houseService)

	flatRepository := flat.NewFlatRepository(db)
	flatService := flat.NewFlatService(flatRepository, publisher)
	flatHandler := flat.NewFlatHandler(flatService)

	router := routes.SetupRouter(authHandler, []byte(cfg.SecretKey), houseHandler, flatHandler)

	mailSender := sender.New()
	subscriber := notification.NewSubscriber(ch, *mailSender, houseRepository)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go subscriber.StartConsuming(ctx)

	fmt.Println("Server is running on port ", cfg.Port)
	if err := http.ListenAndServe(cfg.Port, router); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}

}
