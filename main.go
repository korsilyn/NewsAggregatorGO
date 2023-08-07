package main

import (
	"net/http"
	"net/url"
	"html/template"
	"math"
	"log"
	"os"
	"time"
	"strconv"
	"bytes"

	"github.com/joho/godotenv"
	"news-aggregator-go/news"
)

type SearchQuery struct {
	Query string
	NextPage int
	TotalPages int
	Result *news.Result
}

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

	newsClient := &http.Client{Timeout: 10 * time.Second}
	newsApi := news.NewClient(newsClient, apikey, 20)
	
	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/search", searchHandler(newsApi))

	log.Println("Server started on 127.0.0.1:" + port) 
	http.ListenAndServe(":" + port, mux)
}

func indexHandler (response http.ResponseWriter, request *http.Request) {
	buffer := &bytes.Buffer{}
	err := indexTemplate.Execute(buffer, nil)
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}

	buffer.WriteTo(response)

}

func searchHandler (newsApi *news.Client) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
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
		
		result, err := newsApi.FetchNews(searchQuery, pageQuery)
		
		if err != nil {
			http.Error(response, err.Error(), http.StatusInternalServerError)
			return
		}
		
		nextPage, err := strconv.Atoi(pageQuery)
		if err != nil {
			http.Error(response, err.Error(), http.StatusInternalServerError)
			return
		}

		searchResult := &SearchQuery {
			Query: searchQuery,
			NextPage: nextPage,
			TotalPages: int(math.Ceil(float64(result.TotalResults) / float64(newsApi.MaxResult))),
			Result: result,
		}

		buffer := &bytes.Buffer{}
		err = indexTemplate.Execute(buffer, searchResult)
		if err != nil {
			http.Error(response, err.Error(), http.StatusInternalServerError)
			return
		}

		buffer.WriteTo(response)
		//indexTemplate.Execute(response, nil)
	}
}
