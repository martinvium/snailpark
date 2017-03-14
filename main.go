package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	log.SetFlags(log.Lshortfile)

	// websocket server
	server := NewServer()
	go server.Listen()

	// static files
	http.Handle("/", http.FileServer(http.Dir("public")))

	log.Println("Listening on", ":"+os.Getenv("PORT"))
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), nil))
}
