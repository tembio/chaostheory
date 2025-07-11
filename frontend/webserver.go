package main

import (
	"html/template"
	"log"
	"net/http"
)

var pageTemplate = template.Must(template.ParseFiles("index.html"))

func main() {
	fs := http.FileServer(http.Dir("."))
	http.Handle("/", fs)
	log.Println("Frontend webserver serving static files on :8081")
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
