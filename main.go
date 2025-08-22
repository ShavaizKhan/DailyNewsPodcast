package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	// Simple HTTP client to call Hugging Face API
	"github.com/go-resty/resty/v2"
	// Will extract main content from HTML (removes clutter like buttons, sidebar, etc.)
	"github.com/go-shiori/go-readability"
	// Read .env file
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	// Article URL
	urlStr := os.Getenv("ARTICLE_URL")

	// Fetch HTML
	html, err := fetchHTML(urlStr)
	if err != nil {
		panic(err)
	}

	// Convert string URL to *url.URL
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		panic(err)
	}

	// Extract main content
	article, err := readability.FromReader(strings.NewReader(html), parsedURL)
	if err != nil {
		panic(err)
	}

	fmt.Println("Title:", article.Title)
	fmt.Println("Extracted content length:", len(article.TextContent))

	// Summarize with Hugging Face
	summary, err := summarizeHuggingFace(article.TextContent)
	if err != nil {
		panic(err)
	}

	fmt.Println("Summary:\n", summary)

	// Save summary
	err = os.WriteFile("summary.txt", []byte(summary), 0644)
	if err != nil {
		panic(err)
	}

	fmt.Println(" Summary saved to summary.txt")
}

// fetchHTML fetches the raw HTML of the article
func fetchHTML(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// summarizeHuggingFace sends text to Hugging Face Inference API
func summarizeHuggingFace(text string) (string, error) {
	client := resty.New()
	hfToken := os.Getenv("HF_KEY") // Hugging Face API token
	// Endpoint for summarization model (bart-large-cnn)
	url := os.Getenv("HF_API")

	resp, err := client.R().
		SetHeader("Authorization", "Bearer "+hfToken).
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]string{"inputs": text}).
		Post(url)
	if err != nil {
		return "", err
	}

	// Hugging Face returns JSON array:
	var result []map[string]string
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return "", err
	}

	if len(result) > 0 {
		return result[0]["summary_text"], nil
	}

	return "", fmt.Errorf("no summary returned")
}
