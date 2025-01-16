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
	"os"
	"os/signal"
	"syscall"
	"time"
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

	srv := &http.Server{
		Addr:    cfg.Port,
		Handler: router,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	fmt.Println("Server is running on port", cfg.Port)
	<-quit
	log.Println("Shutting down server...")

	ctx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")

}
