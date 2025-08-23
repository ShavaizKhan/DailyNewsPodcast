package utils

import (
	"bufio"
	"context"
	"io"
	"os"
	"strings"

	//AWS SDKs
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/polly"
	"github.com/aws/aws-sdk-go-v2/service/polly/types"
)

func SynthesizePodcast(script, outputFile string) error {
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
			voice = types.VoiceIdDanielle
			line = strings.TrimPrefix(line, "Alice:")
		} else if strings.HasPrefix(line, "Bob:") {
			voice = types.VoiceIdStephen
			line = strings.TrimPrefix(line, "Bob:")
		} else {
			voice = types.VoiceIdDanielle
		}

		line = strings.TrimSpace(line)

		input := &polly.SynthesizeSpeechInput{
			Text:         &line,
			OutputFormat: types.OutputFormatMp3,
			VoiceId:      voice,
			Engine:       types.EngineGenerative,
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
