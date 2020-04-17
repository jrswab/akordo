package plugins

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
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
func Rule34(req []string, s *dg.Session, msg *dg.MessageCreate) {
	// make sure the channel is marked NSFW
	dChan, err := s.Channel(msg.ChannelID)
	if err != nil {
		log.Printf("session.Channelfailed: %s", err)
	}

	if !dChan.NSFW {
		return
	}

	// Check for proper formatting of message:
	if len(req) < 2 {
		_, err = s.ChannelMessageSend(msg.ChannelID, "Usage: `--rule34 tag`")
		if err != nil {
			log.Printf("session.ChannelMessageSend failed: %s", err)
		}
		return
	}

	// Retrieve an rule34 image based on tag input
	sampleURL, err := requestData(req[1])
	if err != nil {
		log.Printf("failed to request data: %s", err)
	}

	_, err = s.ChannelMessageSend(msg.ChannelID, sampleURL)
	if err != nil {
		log.Printf("session.ChannelMessageSend failed: %s", err)
		return
	}
	log.Printf("%v fetched %s", msg.Member, sampleURL)
}

func requestData(tag string) (string, error) {

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

	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(len(v.Post) - 1)

	return v.Post[n].SampleURL, nil
}
