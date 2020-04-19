package controller

import (
	"fmt"
	"regexp"
	"strings"

	plugs "git.sr.ht/~jrswab/akordo/plugins"
	dg "github.com/bwmarrin/discordgo"
)

const prefix string = `~`

// Controller is used for testing and is implemented by methods that controller how the user
// message gets distributed to the plugins.
type Controller interface {
	CheckSyntax(s *dg.Session, msg *dg.MessageCreate)
	ExecuteTask(req []string, s *dg.Session, msg *dg.MessageCreate)
}

// SessionData holds the data needed for the bot to send/receive messages.
type SessionData struct {
	session *dg.Session

	gifRecord  *plugs.Record
	memeRecord *plugs.Record
	pingRecord *plugs.Record
	r34Record  *plugs.Record
}

// NewSessionData creates a SessionData
func NewSessionData(s *dg.Session) *SessionData {
	return &SessionData{
		session: s,

		gifRecord:  plugs.NewRecorder(),
		memeRecord: plugs.NewRecorder(),
		pingRecord: plugs.NewRecorder(),
		r34Record:  plugs.NewRecorder(),
	}
}

// CheckSyntax uses regexp from the standard library to check the message has the correct
// prefix as defined by the `prefix` constant.
func (sd *SessionData) CheckSyntax(s *dg.Session, msg *dg.MessageCreate) {
	// Make sure the message matches the bot syntax
	regEx := fmt.Sprintf("(?m)^%s(\\w|\\s)+", prefix)
	var re = regexp.MustCompile(regEx)
	match := re.MatchString(msg.Content)
	if !match {
		return
	}

	// Split the string to a slice to parse parameters
	req := strings.Split(msg.Content, " ")
	sd.ExecuteTask(req, s, msg)
}

// ExecuteTask looks up the command found by the bot and kicks off a Goroutine do what
// the user is asking to do.
//
// To remove a plugin simply remove the case statement for that plugin
// To add a plugin, create a case statement for the plugin as shown below.
// If the plugin is new create a new `.go` file under the `plugins` directory.
func (sd *SessionData) ExecuteTask(req []string, s *dg.Session, msg *dg.MessageCreate) {
	switch req[0] {
	case prefix + "gif":
		go sd.gifRecord.Gif(req, s, msg)
	case prefix + "man":
		go plugs.Manual(req, s, msg)
	case prefix + "meme":
		go sd.memeRecord.RequestMeme(req, s, msg)
	case prefix + "ping":
		go sd.pingRecord.Pong(s, msg)
	case prefix + "rule34":
		go sd.r34Record.Rule34(req, s, msg)
	}
}
