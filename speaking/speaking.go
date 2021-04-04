package speaking

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	texttospeech "cloud.google.com/go/texttospeech/apiv1"
	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto"
	texttospeechpb "google.golang.org/genproto/googleapis/cloud/texttospeech/v1"
)

var (
	c *oto.Context
)

func TestAuth() {
	gac := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	if gac == "" {
		panic("GOOGLE_APPLICATION_CREDENTIALS env var not found!")
	}
}

// Say says that a person has left or arrived, returning whether the audio file was from cache
func Say(spokenName string, status string) (usedCache bool, err error) {
	// Check cache
	usedCache = true
	if _, err := os.Stat(fmt.Sprintf("./cachedAudio/%s-%s.mp3", spokenName, status)); os.IsNotExist(err) {
		fmt.Println("File cache not used")
		usedCache = false
		// Make request to GCP and get audio file
		ctx := context.Background()
		client, err := texttospeech.NewClient(ctx)
		if err != nil {
			panic(err)
		}
		req := texttospeechpb.SynthesizeSpeechRequest{
			Input: &texttospeechpb.SynthesisInput{
				InputSource: &texttospeechpb.SynthesisInput_Text{
					Text: fmt.Sprintf("%s has %s", spokenName, status),
				},
			},
			Voice: &texttospeechpb.VoiceSelectionParams{
				LanguageCode: "en-US",
				Name:         "en-US-Wavenet-D",
			},
			AudioConfig: &texttospeechpb.AudioConfig{
				AudioEncoding: texttospeechpb.AudioEncoding_MP3,
				SpeakingRate:  1,
				Pitch:         0,
			},
		}
		resp, err := client.SynthesizeSpeech(ctx, &req)
		if err != nil {
			panic(err)
		}
		// Store file in cache
		err = ioutil.WriteFile(fmt.Sprintf("./cachedAudio/%s-%s.mp3", spokenName, status), resp.AudioContent, 0666)
		if err != nil {
			panic(err)
		}
		fmt.Println("Audio file cached")
	} else {
		fmt.Println("File cache used")
	}
	// Play mp3 file
	f, err := os.Open(fmt.Sprintf("./cachedAudio/%s-%s.mp3", spokenName, status))
	if err != nil {
		panic(err)
	}

	d, err := mp3.NewDecoder(f)
	if err != nil {
		panic(err)
	}

	if c == nil {
		c, err = oto.NewContext(d.SampleRate(), 2, 2, 8192)
		if err != nil {
			panic(err)
		}
	}

	p := c.NewPlayer()
	if _, err := io.Copy(p, d); err != nil {
		return false, err
	}
	fmt.Println(2)

	p.Close()
	f.Close()
	return usedCache, nil
}
