package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	database "server/internal/db"
	// "server/internal/models"
	"server/internal/routes"
	"server/internal/webSocket"

	"github.com/rs/cors"
)

func main() {

	database.ConnectDB()

	// err := database.DB.AutoMigrate(&models.User{}) //uncomment for migration 
	// if err != nil {
	// 	log.Fatalf("Failed to migrate database: %v", err)
	// }

	fmt.Println("Database migration completed!")

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello world using golang")
	})

	http.HandleFunc("/ws", webSocket.HandleConnections)

	router := routes.SetupRoutes()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	handler := c.Handler(router)

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal("Failed to start server:", err)
	}

	fmt.Printf("Server running on port %s\n", port)
	if err := http.Serve(listener, handler); err != nil {
		log.Fatal("Server error:", err)
	}

	log.Fatal(http.ListenAndServe(port, handler))


}
