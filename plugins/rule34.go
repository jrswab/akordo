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

// AkSession allows for tests to mock the discordgo session.Channel() method call
type AkSession interface {
	Channel(channelID string) (st *dg.Channel, err error)
}

type rule34XML struct {
	Count string `xml:"count,attr"`
	Post  []struct {
		SampleURL string `xml:"sample_url,attr"`
	} `xml:"post"`
}

// Rule34Request contains the data to be passed when executing the Rule34() method.
type Rule34Request struct {
	record  *Record
	baseURL string
}

// NewRule34Request creates Rule34Request struct for calling the Rule34 method
// URL is optional to pass in the case the bot maintainer wants to use a different
// gif websites. If more than one URL is passed in the only first will be used.
func NewRule34Request(url ...string) *Rule34Request {
	recorder := NewRecorder()
	rq := &Rule34Request{
		record:  recorder,
		baseURL: fmt.Sprintf("https://rule34.xxx/index.php?page=dapi&s=post&q=index&tags="),
	}

	if len(url) != 0 {
		rq.baseURL = url[0]
	}

	return rq
}

// Rule34 checks that the channel ID is marked as NSFW, makes sure the length of the slice
// is greater that 1 (ie; a tag has been passed with the request) and then retrieves the data.
func (rr *Rule34Request) Rule34(req []string, s AkSession, msg *dg.MessageCreate) (string, error) {
	r := rr.record
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
		return fmt.Sprintf("Usage: `<prefix>rule34 tag`"), nil
	}

	// Check the last time the user made this request
	alertUser, tooSoon := r.CheckLastAsk(msg)
	if tooSoon {
		return alertUser, nil
	}

	url := fmt.Sprintf("%s%s", rr.baseURL, req[1])
	// Retrieve a rule34 image
	sampleURL, err := requestPron(url)
	if err != nil {
		return "", fmt.Errorf("failed to request data: %s", err)
	}

	return sampleURL, nil
}

func requestPron(url string) (string, error) {
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
	returnLen := len(v.Post) - 1

	if returnLen == 0 {
		return v.Post[0].SampleURL, nil
	}

	n := rand.Intn(returnLen)

	return v.Post[n].SampleURL, nil
}
