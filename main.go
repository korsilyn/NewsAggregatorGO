package main

import (
	"net/http"
	"net/url"
	"html/template"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/korsilyn/NewsAggregatorGO/news"
)

var indexTemplate = template.Must(template.ParseFiles("index.html"))

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file!")
	}
	mux := http.NewServeMux()
	port := "3000"
	apikey := os.Getenv("APIKEY")
	if apikey == "" {
		log.Println("No api key!")
	}

	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/search", searchHandler)

	log.Println("Server started on 127.0.0.1:" + port) 
	http.ListenAndServe(":" + port, mux)
}

func indexHandler (response http.ResponseWriter, request *http.Request) {
	indexTemplate.Execute(response, nil)
}

func searchHandler (response http.ResponseWriter, request *http.Request) {
	url_var, err := url.Parse(request.URL.String())
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}

	params := url_var.Query()
	searchQuery := params.Get("keyword")
	pageQuery := params.Get("page")
	if pageQuery == "" {
		pageQuery = "1"
	}

	fmt.Println("Search query: ", searchQuery)
	fmt.Println("Page: ", pageQuery)

	indexTemplate.Execute(response, nil)
}
