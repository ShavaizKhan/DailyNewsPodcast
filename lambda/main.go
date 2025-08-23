package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/ShavaizKhan/DailyNewsPodcast/utils"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Event struct{}

type Response struct {
	StatusCode int    `json:"statusCode"`
	Body       string `json:"body"`
}

func HandleRequest(ctx context.Context, event Event) (Response, error) {
	newsAPIKey := os.Getenv("NEWS_KEY")
	groqToken := os.Getenv("GROQ_KEY")
	s3Bucket := os.Getenv("S3_BUCKET")

	// Get news articles
	articles, err := utils.FetchNews(newsAPIKey)
	if err != nil {
		return Response{
			StatusCode: 500,
			Body:       fmt.Sprintf("Error fetching news: %v", err),
		}, err
	}

	// Generate podcast script
	var podcastScript string
	for i, article := range articles {
		dialogue, err := utils.GenerateDialogue(article, groqToken)
		if err != nil {
			return Response{
				StatusCode: 500,
				Body:       fmt.Sprintf("Error generating dialogue: %v", err),
			}, err
		}

		if i == 0 {
			podcastScript += fmt.Sprintf("\nWelcome back to your daily news update!\n%s", dialogue)
		} else if i < len(articles)-1 {
			podcastScript += fmt.Sprintf("\nMoving on to our next discussion.\n%s", dialogue)
		} else {
			podcastScript += fmt.Sprintf("\nNow to our final story.\n%s", dialogue)
		}
	}
	podcastScript += "\n\nThank you for tuning in! We'll be back with more news coverage for you tomorrow!"

	// Create temporary file for audio
	fileName := fmt.Sprintf("general_podcast_%s.mp3", time.Now().Format("2006-01-02"))
	tmpFile := "/tmp/" + fileName

	// Generate audio
	err = utils.SynthesizePodcast(podcastScript, tmpFile)
	if err != nil {
		return Response{
			StatusCode: 500,
			Body:       fmt.Sprintf("Error synthesizing podcast: %v", err),
		}, err
	}

	// Upload to S3
	err = uploadToS3(tmpFile, s3Bucket, fileName)
	if err != nil {
		return Response{
			StatusCode: 500,
			Body:       fmt.Sprintf("Error uploading to S3: %v", err),
		}, err
	}

	// Clean up temp file
	os.Remove(tmpFile)

	return Response{
		StatusCode: 200,
		Body:       fmt.Sprintf("Podcast generated successfully: %s", fileName),
	}, nil
}

func uploadToS3(filePath, bucket, key string) error {
	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return err
	}

	client := s3.NewFromConfig(cfg)

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: &bucket,
		Key:    &key,
		Body:   file,
	})

	return err
}

func main() {
	lambda.Start(HandleRequest)
}
