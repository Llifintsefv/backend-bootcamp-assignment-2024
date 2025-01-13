package routes

import (
	"backend-bootcamp-assignment-2024/internal/auth"
	"backend-bootcamp-assignment-2024/internal/flat"
	"backend-bootcamp-assignment-2024/internal/house"

	"github.com/gorilla/mux"
)

func SetupRouter(authHandler *auth.AuthHandler, secretKey []byte, houseHandler *house.HouseHandler, flatHandler  *flat.FlatHandler) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/dummyLogin", authHandler.DummyLoginHandler).Methods("GET")
	router.HandleFunc("/login", authHandler.Login).Methods("POST")
	router.HandleFunc("/register", authHandler.Register).Methods("POST")
	router.HandleFunc("/flat/{id}",flatHandler.GetFlats).Methods("GET")

	authRouter := router.PathPrefix("/").Subrouter()
	authRouter.Use(auth.AuthMiddleware(secretKey))
	authRouter.HandleFunc("/flat/create",flatHandler.CreateFlat).Methods("POST")
	



	moderatorRouter := router.PathPrefix("/").Subrouter()
	moderatorRouter.Use(auth.AuthMiddleware(secretKey), auth.ModeratorMiddleware)
	moderatorRouter.HandleFunc("/house/create", houseHandler.CreateHouse).Methods("POST")
	moderatorRouter.HandleFunc("/flat/update",flatHandler.UpdateFlatStatus).Methods("POST")

	return router
}
