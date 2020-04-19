package plugins

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"

	dg "github.com/bwmarrin/discordgo"
)

type rule34XML struct {
	Count string `xml:"count,attr"`
	Post  []struct {
		SampleURL string `xml:"sample_url,attr"`
	} `xml:"post"`
}

// Rule34 checks that the channel ID is marked as NSFW, makes sure the length of the slice
// is greater that 1 (ie; a tag has been passed with the request) and then retrieves the data.
func (r *Record) Rule34(req []string, s *dg.Session, msg *dg.MessageCreate) (string, error) {
	// make sure the channel is marked NSFW
	dChan, err := s.Channel(msg.ChannelID)
	if err != nil {
		return "", fmt.Errorf("session.Channelfailed: %s", err)
	}

	if !dChan.NSFW {
		return "", nil
	}

	// Check for proper formatting of message:
	if len(req) < 2 {
		return fmt.Sprintf("Usage: `--rule34 tag`"), nil
	}

	// Check the last time the user made this request
	alertUser, tooSoon := r.checkLastAsk(s, msg)
	if tooSoon {
		return alertUser, nil
	}

	// Retrieve an rule34 image based on tag input
	sampleURL, err := requestPron(req[1])
	if err != nil {
		return "", fmt.Errorf("failed to request data: %s", err)
	}

	return sampleURL, nil
}

func requestPron(tag string) (string, error) {

	url := fmt.Sprintf("https://rule34.xxx/index.php?page=dapi&s=post&q=index&tags=%s", tag)
	res, err := http.Get(url)
	if err != nil {
		return "", err
	}

	xmlData, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return "", err
	}

	v := &rule34XML{}
	err = xml.Unmarshal([]byte(xmlData), &v)
	if err != nil {
		return "", err
	}

	// If look up returns an empty slice display this message instead.
	if len(v.Post) < 1 {
		return "No results found :sob:", nil
	}
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(len(v.Post) - 1)

	return v.Post[n].SampleURL, nil
}
