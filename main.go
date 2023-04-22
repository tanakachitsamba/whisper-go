package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/joho/godotenv"
	openai "github.com/sashabaranov/go-openai"
)

func transcribeFile(c *openai.Client, ctx context.Context, audioFile string, wg *sync.WaitGroup, transcriptions *[]string) {
	defer wg.Done()

	req := openai.AudioRequest{
		Model:    openai.Whisper1,
		FilePath: audioFile,
	}
	resp, err := c.CreateTranscription(ctx, req)
	if err != nil {
		fmt.Printf("Transcription error for %s: %v\n", audioFile, err)
		return
	}

	*transcriptions = append(*transcriptions, resp.Text)
}

func whisper(c *openai.Client, ctx context.Context) {

	// Define the audio files directory and the audio format
	audioDir := "audios/"
	audioFormat := ".mp3"

	// Read the audio files from the directory
	files, err := ioutil.ReadDir(audioDir)
	if err != nil {
		fmt.Printf("Error reading audio directory: %v\n", err)
		return
	}

	// Filter audio files based on the audio format
	var audioFiles []string
	for _, file := range files {
		if strings.HasSuffix(file.Name(), audioFormat) {
			audioFiles = append(audioFiles, audioDir+file.Name())
		}
	}

	// Process the audio files and store their transcriptions in a slice
	var transcriptions []string
	var wg sync.WaitGroup
	for _, audioFile := range audioFiles {
		wg.Add(1)
		go transcribeFile(c, ctx, audioFile, &wg, &transcriptions)
	}

	wg.Wait()

	// Print the transcriptions
	for i, transcription := range transcriptions {
		fmt.Printf("Transcription %d: %s\n", i+1, transcription)
	}

}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("error loading .env file:", err)
	}
	c := openai.NewClient(os.Getenv("OPENAI_KEY"))

	whisper(c, context.Background())
}
