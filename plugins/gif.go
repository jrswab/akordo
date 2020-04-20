package plugins

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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

// GifRequest contains the data to be passed when executing the Gif() method.
type GifRequest struct {
	Record  *Record
	BaseURL string
}

// NewGifRequest creates GifRequest struct for calling the Gif method
// URL is optional to pass in the case the bot maintainer wants to use a different
// gif websites. If more than one URL is passed in the only first will be used.
func NewGifRequest(url ...string) *GifRequest {
	recorder := NewRecorder()
	rq := &GifRequest{
		Record:  recorder,
		BaseURL: fmt.Sprintf("http://api.giphy.com/v1/gifs/search?rating=pg"),
	}

	if len(url) != 0 {
		rq.BaseURL = url[0]
	}

	return rq
}

// Gif makes sure the length of the slice is greater that 1
// (ie; a tag has been passed with the request) and then return a random gif from Giphy.
func (g *GifRequest) Gif(req []string, s *dg.Session, msg *dg.MessageCreate) (string, error) {
	r := g.Record
	// Check the last time the user made this request
	alertUser, tooSoon := r.checkLastAsk(msg)
	if tooSoon {
		return alertUser, nil
	}

	// Check for proper formatting of message:
	if len(req) < 2 {
		return fmt.Sprintf("Usage: `<prefix>gif word`"), nil
	}

	// Create URL
	giphyAPI := os.Getenv("GIPHY_KEY")
	url := fmt.Sprintf("%s&api_key=%s&q=%s", g.BaseURL, giphyAPI, req[1])

	// Retrieve an rule34 image based on tag input
	sampleURL, err := requestGif(url)
	if err != nil {
		return "", fmt.Errorf("failed to request gif data: %s", err)
	}

	return sampleURL, nil
}

func requestGif(url string) (string, error) {
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
