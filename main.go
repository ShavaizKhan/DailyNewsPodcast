package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	// Simple HTTP client to call Hugging Face API
	"github.com/go-resty/resty/v2"
	// Read .env file
	"github.com/joho/godotenv"
	// AWS SDKs

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/polly"
	"github.com/aws/aws-sdk-go-v2/service/polly/types"
)

type NewsAPIResponse struct {
	Articles []struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Content     string `json:"content"`
		URL         string `json:"url"`
	} `json:"articles"`
}
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

func fetchTopHeadlines(apiKey string) ([]string, []string, error) {
	url := fmt.Sprintf("https://newsapi.org/v2/top-headlines?country=us&pageSize=10&apiKey=%s", apiKey)

	resp, err := http.Get(url)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	var data NewsAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, nil, err
	}

	var articles, titles []string
	for _, a := range data.Articles {
		articles = append(articles, fmt.Sprintf("%s - %s", a.Title, a.Description))
		titles = append(titles, a.Title)
	}
	return titles, articles, nil
}

func generateDialogue(article, groqToken string) (string, error) {
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

func synthesizePodcast(script, outputFile string) error {
	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1"))
	if err != nil {
		return err
	}

	client := polly.NewFromConfig(cfg)

	outFile, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer outFile.Close()

	scanner := bufio.NewScanner(strings.NewReader(script))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		// Determine speaker voice and remove prefix
		var voice types.VoiceId
		if strings.HasPrefix(line, "Alice:") {
			voice = types.VoiceIdRuth
			line = strings.TrimPrefix(line, "Alice:")
		} else if strings.HasPrefix(line, "Bob:") {
			voice = types.VoiceIdStephen
			line = strings.TrimPrefix(line, "Bob:")
		} else {
			voice = types.VoiceIdRuth
		}

		line = strings.TrimSpace(line)

		input := &polly.SynthesizeSpeechInput{
			Text:         &line,
			OutputFormat: types.OutputFormatMp3,
			VoiceId:      voice,
			Engine:       types.EngineNeural,
		}

		resp, err := client.SynthesizeSpeech(ctx, input)
		if err != nil {
			return err
		}
		defer resp.AudioStream.Close()

		if _, err := io.Copy(outFile, resp.AudioStream); err != nil {
			return err
		}
	}

	return scanner.Err()
}

func main() {
	godotenv.Load()
	newsAPIKey := os.Getenv("NEWS_KEY")
	groqToken := os.Getenv("GROQ_KEY")

	// Get news articles
	titles, articles, err := fetchTopHeadlines(newsAPIKey)
	if err != nil {
		panic(err)
	}

	// Generate podcast-style dialogue
	var podcastScript string
	for i, article := range articles {
		dialogue, err := generateDialogue(article, groqToken)
		if err != nil {
			panic(err)
		}
		if 1 <= i && i < len(articles)-1 {
			podcastScript += fmt.Sprintf("\nMoving on to our next discussion - %s\n%s", titles[i], dialogue)
		} else if i < 1 {
			podcastScript += fmt.Sprintf("\nWelcome back to your daily news update!\n%s", dialogue)
		} else {
			podcastScript += fmt.Sprintf("\nNow to our final story.\n%s", dialogue)
		}
	}
	podcastScript += "\n\nThank you for tuning in! We'll be back with more news coverage for you tomorrow!"
	// fmt.Println("Generated podcast script:\n", podcastScript)

	e := synthesizePodcast(podcastScript, fmt.Sprintf("podcasts/general_podcast_%s.mp3", time.Now().Format("2006-01-02")))
	if e != nil {
		fmt.Println("Error generating podcast:", e)
	} else {
		fmt.Println("Podcast saved!")
	}
}
