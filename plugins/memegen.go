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

// RequestMeme receives the users request for a meme with the given parameters.
// If the resquest is malformed (ie, only one word after --meme) the function
// terminates and returns a message to the sure on how to use the meme generator.
func (r *Record) RequestMeme(req []string, s *dg.Session, msg *dg.MessageCreate) (string, error) {
	if len(req) < 3 {
		switch req[1] {
		case "list":
			listMsg0 := "To see all available memes head to https://memegen.link/api/templates/\n"
			listMsg1 := "Use the name at the end of the URLs that are displayed."
			return (listMsg0 + listMsg1), nil
		default:
			return fmt.Sprintf("Usage: `--meme name top_text <bottom_text>`"), nil
		}
	}

	if len(req) > 4 {
		return fmt.Sprintf("Usage: `--meme name top_text <bottom_text>`"), nil
	}

	// Check the last time the user made this request
	alertUser, tooSoon := r.checkLastAsk(msg)
	if tooSoon {
		return alertUser, nil
	}

	// Retrieve the generated meme based on tag input
	URL, err := generateMeme(req)
	if err != nil {
		return "", fmt.Errorf("generateMeme failed: %s", err)
	}

	return URL, nil
}

func generateMeme(req []string) (string, error) {
	url := fmt.Sprintf("https://memegen.link/api/templates/")
	for idx, word := range req {
		if idx == 0 {
			continue
		}

		url = fmt.Sprintf("%s%s/", url, word)
	}

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
