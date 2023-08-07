package main

import (
	"net/http"
	"html/template"
	"log"
)

var indexTemplate = template.Must(template.ParseFiles("index.html"))

func indexHandler (response http.ResponseWriter, request *http.Request) {
	indexTemplate.Execute(response, nil)
}

func main() {
	mux := http.NewServeMux()
	port := "3000"

	mux.HandleFunc("/", indexHandler)
	log.Println("Server started on 127.0.0.1:" + port) 
	http.ListenAndServe(":" + port, mux)
}
