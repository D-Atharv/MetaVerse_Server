package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"server/internal/db"
	"server/internal/routes"
	"server/internal/webSocket"

	"github.com/rs/cors"
)

func main() {
	database.ConnectDB()

	// Uncomment for database migration (when needed)
	// err := db.DB.AutoMigrate(&models.User{})
	// if err != nil {
	// 	log.Fatalf("Database migration failed: %v", err)
	// }
	// fmt.Println("Database migration completed!")

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "http://localhost:3000"}, 
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	mux := http.NewServeMux()

	mux.HandleFunc("/ws", webSocket.HandleConnections)

	router := routes.SetupRoutes()
	mux.Handle("/", router)

	handler := c.Handler(mux)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" 
	}

	serverAddress := fmt.Sprintf(":%s", port)
	fmt.Printf("Server running on port %s\n", port)

	if err := http.ListenAndServe(serverAddress, handler); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
