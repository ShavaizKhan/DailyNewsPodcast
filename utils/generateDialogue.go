package utils

import (
	"fmt"
	// Simple HTTP client to call API
	"github.com/go-resty/resty/v2"
)

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Model       string        `json:"model"`
	Messages    []ChatMessage `json:"messages"`
	Stream      bool          `json:"stream"`
	MaxTokens   int           `json:"max_tokens,omitempty"`
	Temperature float64       `json:"temperature,omitempty"`
}

type ChatChoice struct {
	Message struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"message"`
}

type ChatResponse struct {
	Choices []ChatChoice `json:"choices"`
}

func GenerateDialogue(article, groqToken string) (string, error) {
	prompt := fmt.Sprintf("Turn this article into a short podcast-style conversation between two hosts, Alice and Bob without any intro and outro. Keep it engaging but concise, and sounding natural. Keep it within 1000 characters and make a new line for each speaker with the prefix 'Bob:' or 'Alice:'. Ensure there's a newline between each speaker :\n\n%s", article)

	client := resty.New()

	request := ChatRequest{
		Model: "llama-3.1-8b-instant",
		Messages: []ChatMessage{
			{
				Role:    "user",
				Content: prompt,
			},
		},
		MaxTokens:   400,
		Temperature: 0.7,
	}

	var result ChatResponse
	resp, err := client.R().
		SetHeader("Authorization", "Bearer "+groqToken).
		SetHeader("Content-Type", "application/json").
		SetBody(request).
		SetResult(&result). // This automatically unmarshals successful responses
		Post("https://api.groq.com/openai/v1/chat/completions")

	if err != nil {
		return "", err
	}

	if resp.StatusCode() != 200 {
		return "", fmt.Errorf("groq api error (status %d): %s", resp.StatusCode(), resp.String())
	}

	if len(result.Choices) > 0 {
		return result.Choices[0].Message.Content, nil
	}

	return "", fmt.Errorf("no response generated")
}
