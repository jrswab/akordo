package plugins

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	dg "github.com/bwmarrin/discordgo"
)

type giphyData struct {
	Data []struct {
		EmbedURL string `json:"embed_url"`
	} `json:"data"`
}

// Gif makes sure the length of the slice is greater that 1
// (ie; a tag has been passed with the request) and then return a random gif from Giphy.
func (r *Record) Gif(req []string, s *dg.Session, msg *dg.MessageCreate) {
	// Check the last time the user made this request
	if tooSoon := r.checkLastAsk(s, msg); tooSoon {
		return
	}

	// Check for proper formatting of message:
	if len(req) < 2 {
		_, err := s.ChannelMessageSend(msg.ChannelID, "Usage: `--gif word`")
		if err != nil {
			log.Printf("session.ChannelMessageSend failed: %s", err)
		}
		return
	}

	// Retrieve an rule34 image based on tag input
	sampleURL, err := requestGif(req[1])
	if err != nil {
		log.Printf("failed to request data: %s", err)
	}

	_, err = s.ChannelMessageSend(msg.ChannelID, sampleURL)
	if err != nil {
		log.Printf("session.ChannelMessageSend failed: %s", err)
		return
	}
	log.Printf("%s fetched gif: %s", msg.Author.Username, sampleURL)
}

func requestGif(tag string) (string, error) {
	giphyAPI := os.Getenv("GIPHY_KEY")

	url := fmt.Sprintf("http://api.giphy.com/v1/gifs/search?api_key=%s&rating=pg&q=%s", giphyAPI, tag)
	res, err := http.Get(url)
	if err != nil {
		return "", err
	}

	gifData, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return "", err
	}

	v := &giphyData{}
	err = json.Unmarshal([]byte(gifData), &v)
	if err != nil {
		return "", err
	}

	// If look up returns an empty slice display this message instead.
	if len(v.Data) < 1 {
		return "No results found :sob:", nil
	}
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(len(v.Data) - 1)

	return v.Data[n].EmbedURL, nil
}
