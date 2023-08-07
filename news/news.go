package news

import "net/http"

type Client struct {
	httpClient *http.Client
	key string
	maxResult int
}

func NewClient(httpClient *http.Client, key string, maxResult int) {
	if maxResult >= 20 {
		maxResult = 20
	}

	return &Client(httpClient, key, maxResult)
}
