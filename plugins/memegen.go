package plugins

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	dg "github.com/bwmarrin/discordgo"
)

type memeGenJSON struct {
	Direct struct {
		Masked string `json:"masked"`
	} `json:"direct"`
}

// MemeRequest contains the data to be passed when executing the Gif() method.
type MemeRequest struct {
	record  *Record
	baseURL string
}

// NewMemeRequest creates GifRequest struct for calling the Gif method
// URL is optional to pass in the case the bot maintainer wants to use a different
// gif websites. If more than one URL is passed in the only first will be used.
func NewMemeRequest(url ...string) *MemeRequest {
	recorder := NewRecorder()
	rq := &MemeRequest{
		record:  recorder,
		baseURL: fmt.Sprintf("https://memegen.link/api/templates/"),
	}

	if len(url) != 0 {
		rq.baseURL = url[0]
	}

	return rq
}

// RequestMeme receives the users request for a meme with the given parameters.
// If the resquest is malformed (ie, only one word after --meme) the function
// terminates and returns a message to the sure on how to use the meme generator.
func (m *MemeRequest) RequestMeme(req []string, s *dg.Session, msg *dg.MessageCreate) (string, error) {
	usageReturn := fmt.Sprintf("Usage: `<prefix>meme name top_text <bottom_text>`")

	if len(req) < 3 {
		switch req[1] {
		case "list":
			listMsg0 := "To see all available memes head to https://memegen.link/api/templates/\n"
			listMsg1 := "Use the name at the end of the URLs that are displayed."
			return (listMsg0 + listMsg1), nil
		default:
			return usageReturn, nil
		}
	}

	if len(req) > 4 {
		return usageReturn, nil
	}

	// Check the last time the user made this request
	alertUser, tooSoon := m.record.CheckLastAsk(msg)
	if tooSoon {
		return alertUser, nil
	}

	// Format the URL
	url := fmt.Sprintf("%s", m.baseURL)
	for idx, word := range req {
		if idx == 0 {
			continue
		}

		url = fmt.Sprintf("%s%s/", url, word)
	}

	// Retrieve the generated meme based on tag input
	URL, err := generateMeme(url)
	if err != nil {
		return "", fmt.Errorf("generateMeme failed: %s", err)
	}

	return URL, nil
}

func generateMeme(url string) (string, error) {
	res, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to get URL: %s", err)
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("unable to read URL body: %s", err)
	}
	res.Body.Close()

	v := &memeGenJSON{}
	err = json.Unmarshal([]byte(data), &v)
	if err != nil {
		return "", fmt.Errorf("unmarshal or memegen JSON fail.d: %s", err)
	}

	return v.Direct.Masked, nil
}
