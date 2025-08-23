package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type NewsAPIResponse struct {
	Articles []struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Content     string `json:"content"`
		URL         string `json:"url"`
	} `json:"articles"`
}

func FetchNews(apiKey string) ([]string, error) {
	url := fmt.Sprintf("https://newsapi.org/v2/top-headlines?country=us&pageSize=10&apiKey=%s", apiKey)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data NewsAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	var articles []string
	for _, a := range data.Articles {
		articles = append(articles, fmt.Sprintf("%s - %s", a.Title, a.Description))
	}
	return articles, nil
}
