package controller

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"git.sr.ht/~jrswab/akordo/plugins"
	plugs "git.sr.ht/~jrswab/akordo/plugins"
	dg "github.com/bwmarrin/discordgo"
)

// Controller is used for testing and is implemented by methods that controller how the user
// message gets distributed to the plugins.
type Controller interface {
	CheckSyntax(s *dg.Session, msg *dg.MessageCreate)
	ExecuteTask(req []string, s *dg.Session, msg *dg.MessageCreate)
}

// SessionData holds the data needed to complete the requested transactions
type SessionData struct {
	session *dg.Session
	prefix  string

	crypto      *plugins.Crypto
	gifRequest  *plugs.GifRequest
	memeRequest *plugs.MemeRequest
	pingRecord  *plugs.Record
	r34Request  *plugs.Rule34Request
}

// NewSessionData creates a SessionData
func NewSessionData(s *dg.Session) *SessionData {
	return &SessionData{
		session: s,
		prefix:  `=`,

		crypto:      plugs.NewCrypto(),
		gifRequest:  plugs.NewGifRequest(),
		memeRequest: plugs.NewMemeRequest(),
		pingRecord:  plugs.NewRecorder(),
		r34Request:  plugs.NewRule34Request(),
	}
}

// NewMessage waits for a ne message to be sent in a the Discord guild
// This kicks off a Goroutine to free up the mutex set by discordgo `AddHandler` method.
func (sd *SessionData) NewMessage(s *dg.Session, msg *dg.MessageCreate) {
	go sd.checkSyntax(s, msg)
}

// CheckSyntax uses regexp from the standard library to check the message has the correct
// prefix as defined by the `prefix` constant.
func (sd *SessionData) checkSyntax(s *dg.Session, msg *dg.MessageCreate) {
	// Make sure the message matches the bot syntax
	regEx := fmt.Sprintf("(?m)^%s(\\w|\\s)+", sd.prefix)
	var re = regexp.MustCompile(regEx)
	match := re.MatchString(msg.Content)
	if !match {
		return
	}

	sd.ExecuteTask(s, msg)
}

// ExecuteTask looks up the command found by the bot and kicks off a Goroutine do what
// the user is asking to do.
//
// To remove a plugin simply remove the case statement for that plugin
// To add a plugin, create a case statement for the plugin as shown below.
// If the plugin is new create a new `.go` file under the `plugins` directory.
func (sd *SessionData) ExecuteTask(s *dg.Session, msg *dg.MessageCreate) {
	var res string
	var err error
	var isDM bool

	// Split the string to a slice to parse parameters
	req := strings.Split(msg.Content, " ")

	switch req[0] {
	case sd.prefix + "crypto":
		res, err = sd.crypto.Game(req, msg)
	case sd.prefix + "gif":
		res, err = sd.gifRequest.Gif(req, s, msg)
	case sd.prefix + "man":
		isDM = true
		res = plugs.Manual(req, s, msg)
	case sd.prefix + "meme":
		res, err = sd.memeRequest.RequestMeme(req, s, msg)
	case sd.prefix + "ping":
		res = sd.pingRecord.Pong(msg)
	case sd.prefix + "rule34":
		res, err = sd.r34Request.Rule34(req, s, msg)
	}

	sd.Reply(res, err, isDM, msg)
}

// Reply takes the executed data and replies to the user. This is either in the channel
// where the command was sent or as a direct message to the user.
func (sd *SessionData) Reply(res string, err error, isDM bool, msg *dg.MessageCreate) {
	s := sd.session

	if err != nil {
		log.Printf("error executing task: %s", err)
		return
	}

	if res == "" {
		return
	}

	if isDM {
		dm, err := s.UserChannelCreate(msg.Author.ID)
		if err != nil {
			log.Printf("s.UserChannelCreate failed to create DM for %s: %s",
				msg.Author.Username, err)
			return
		}

		_, err = s.ChannelMessageSend(dm.ID, res)
		if err != nil {
			log.Printf("session.ChannelMessageSend failed to send DM: %s", err)
		}
		return
	}

	_, err = s.ChannelMessageSend(msg.ChannelID, res)
	if err != nil {
		log.Printf("session.ChannelMessageSend failed: %s", err)
	}
}
