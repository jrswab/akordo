package plugins

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
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
func (r *Record) RequestMeme(req []string, s *dg.Session, msg *dg.MessageCreate) {
	if len(req) < 3 {
		switch req[1] {
		case "list":
			listMsg0 := "To see all available memes head to https://memegen.link/api/templates/\n"
			listMsg1 := "Use the name at the end of the URLs that are displayed."
			_, err := s.ChannelMessageSend(msg.ChannelID, (listMsg0 + listMsg1))
			if err != nil {
				log.Printf("session.ChannelMessageSend failed: %s", err)
			}
			return
		default:
			_, err := s.ChannelMessageSend(msg.ChannelID, "Usage: `--meme name top_text <bottom_text>`")
			if err != nil {
				log.Printf("session.ChannelMessageSend failed: %s", err)
			}
			return
		}
	}

	if len(req) > 4 {
		_, err := s.ChannelMessageSend(msg.ChannelID, "Usage: `--meme name top_text <bottom_text>`")
		if err != nil {
			log.Printf("session.ChannelMessageSend failed: %s", err)
		}
		return
	}

	// Check the last time the user made this request
	if tooSoon := r.checkLastAsk(s, msg); tooSoon {
		return
	}

	// Retrieve the generated meme based on tag input
	URL, err := generateMeme(req)
	if err != nil {
		log.Printf("generateMeme failed: %s", err)
		return
	}

	_, err = s.ChannelMessageSend(msg.ChannelID, URL)
	if err != nil {
		log.Printf("session.ChannelMessageSend failed: %s", err)
		return
	}
	log.Printf("%s generated meme: %s", msg.Member.User.Username, URL)
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
