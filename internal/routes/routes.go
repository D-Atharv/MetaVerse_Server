package routes

import (
	"net/http"
	"server/internal/auth"

	"github.com/gorilla/mux"
)

func SetupRoutes() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/register", auth.RegisterHandler).Methods("POST")
	router.HandleFunc("/login", auth.LoginHandler).Methods("POST")

	router.HandleFunc("/protected", auth.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to the protected route"))
	})).Methods("GET")

	return router
}
