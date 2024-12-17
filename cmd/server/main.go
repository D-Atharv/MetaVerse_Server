package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	database "server/internal/db"
	"server/internal/routes"
	"server/internal/webSocket"
)

func main() {

	database.ConnectDB()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello world using golang")
	})

	http.HandleFunc("/ws", webSocket.HandleConnections)

	router := routes.SetupRoutes()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal("Failed to start server:", err)
	}

	fmt.Printf("Server running on port %s\n", port)
	if err := http.Serve(listener, nil); err != nil {
		log.Fatal("Server error:", err)
	}

	log.Fatal(http.ListenAndServe(port, router))

}
