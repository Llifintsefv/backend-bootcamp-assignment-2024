package routes

import (
	"backend-bootcamp-assignment-2024/internal/auth"

	"github.com/gorilla/mux"
)

func SetupRouter(authHandler *auth.AuthHandler, secretKey []byte, houseHandler *house.HouseHandler) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/dummyLogin", authHandler.DummyLoginHandler).Methods("GET")
	router.HandleFunc("/login", authHandler.Login).Methods("POST")
	router.HandleFunc("/register", authHandler.Register).Methods("POST")

	authRouter := router.PathPrefix("/").Subrouter()
	authRouter.Use(auth.AuthMiddleware(secretKey))
	authRouter.HandleFunc("/house/{id}", houseHandler.GetHouseHandler).Methods("GET")

	moderatorRouter := router.PathPrefix("/").Subrouter()
	moderatorRouter.Use(auth.AuthMiddleware(secretKey), auth.ModeratorMiddleware)
	moderatorRouter.HandleFunc("/house/create", houseHandler.CreateHouseHandler).Methods("POST")

	return router
}
