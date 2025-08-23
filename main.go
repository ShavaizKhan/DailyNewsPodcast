package main

import (
	"fmt"
	"os"
	"time"

	// Read .env file
	"github.com/joho/godotenv"
	// Util functions
	"github.com/ShavaizKhan/DailyNewsPodcast/utils"
)

func main() {
	godotenv.Load()
	newsAPIKey := os.Getenv("NEWS_KEY")
	groqToken := os.Getenv("GROQ_KEY")

	// Get news articles
	articles, err := utils.FetchNews(newsAPIKey)
	if err != nil {
		panic(err)
	}

	// Generate podcast-style dialogue
	var podcastScript string
	for i, article := range articles {
		dialogue, err := utils.GenerateDialogue(article, groqToken)
		if err != nil {
			panic(err)
		}
		if 1 <= i && i < len(articles)-1 {
			podcastScript += fmt.Sprintf("\nMoving on to our next discussion.\n%s", dialogue)
		} else if i < 1 {
			podcastScript += fmt.Sprintf("\nWelcome back to your daily news update!\n%s", dialogue)
		} else {
			podcastScript += fmt.Sprintf("\nNow to our final story.\n%s", dialogue)
		}
	}
	podcastScript += "\n\nThank you for tuning in! We'll be back with more news coverage for you tomorrow!"
	// fmt.Println("Generated podcast script:\n", podcastScript)

	e := utils.SynthesizePodcast(podcastScript, fmt.Sprintf("general_podcast_%s.mp3", time.Now().Format("2006-01-02")))
	if e != nil {
		fmt.Println("Error generating podcast:", e)
	} else {
		fmt.Println("Podcast saved!")
	}
}
