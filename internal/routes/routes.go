package routes

import (
	"backend-bootcamp-assignment-2024/internal/auth"

	"github.com/gorilla/mux"
)

func SetupRouter(authHandler *auth.AuthHandler) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/dummyLogin", authHandler.DummyLoginHandler).Methods("GET")
	router.HandleFunc("/login", authHandler.Login).Methods("POST")
	router.HandleFunc("/register", authHandler.Register).Methods("POST")

	return router
}
