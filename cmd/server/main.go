package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"server/internal/webSocket"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello world using golang")
	})

	http.HandleFunc("/ws", webSocket.HandleConnections)

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Server running at port 8080")
	err = http.Serve(listener, nil)

	if err != nil {
		log.Fatal(err)
	}
}
