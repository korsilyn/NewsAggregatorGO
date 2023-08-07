package news

import (
	"net/http"
	"net/url"
	"encoding/json"
	"io/ioutil"
	"fmt"
	"time"
)

type Client struct {
	httpClient *http.Client
	key string
	MaxResult int
}

type Result struct {
	Status       string `json:"status"`
	TotalResults int    `json:"totalResults"`
	Articles     []Article `json:"articles"`
}

type Article struct {
	Source struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"source"`
	Author      string    `json:"author"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	URL         string    `json:"url"`
	URLToImage  string    `json:"urlToImage"`
	PublishedAt time.Time `json:"publishedAt"`
	Content     string    `json:"content"`
}

func NewClient(httpClient *http.Client, key string, maxResult int) *Client {
	if maxResult >= 50 {
		maxResult = 50
	}

	return &Client{httpClient, key, maxResult}
}

func (c *Client) FetchNews(query, page string) (*Result, error) {
	address := fmt.Sprintf("https://newsapi.org/v2/everything?q=%s&pageSize=%d&page=%s&apiKey=%s&sortBy=publishedAt&language=en", url.QueryEscape(query), c.MaxResult, page, c.key)
	response, err := c.httpClient.Get(address)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(string(body))
	}

	result := &Result{}
	return result, json.Unmarshal(body, result)
}

func (article *Article) FormatDate() string {
	year, month, day := article.PublishedAt.Date()
	return fmt.Sprintf("%v %d, %d", month, day, year)
}
