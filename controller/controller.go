package controller

import (
	"regexp"
	"strings"

	plugs "git.sr.ht/~jrswab/akordo/plugins"
	dg "github.com/bwmarrin/discordgo"
)

// Controller holds the data needed for the bot to send/receive messages.
type Controller struct {
	session *dg.Session

	gifRecord  *plugs.Record
	memeRecord *plugs.Record
	pingRecord *plugs.Record
	r34Record  *plugs.Record
}

// NewController creates a Controller
func NewController(s *dg.Session) *Controller {
	return &Controller{
		session: s,

		gifRecord:  plugs.NewRecorder(),
		memeRecord: plugs.NewRecorder(),
		pingRecord: plugs.NewRecorder(),
		r34Record:  plugs.NewRecorder(),
	}
}

// ReceiveMessage reads the message and returns the response based on the input.
func (c *Controller) ReceiveMessage(s *dg.Session, msg *dg.MessageCreate) {
	// Make sure the message matches the bot syntax
	var re = regexp.MustCompile(`(?m)^--(\w|\s)+`)
	match := re.MatchString(msg.Content)
	if !match {
		return
	}

	// Split the string to a slice to parse parameters
	req := strings.Split(msg.Content, " ")

	// Perform action based off each command
	// To remove a plugin simply remove the case statement for that plugin
	// To add a plugin, create a case statement for the plugin as shown below.
	// If the plugin is new create a new `.go` file under the `plugins` directory.
	switch req[0] {
	case "--gif":
		go c.gifRecord.Gif(req, s, msg)
	case "--man":
		go plugs.Manual(req, s, msg)
	case "--meme":
		go c.memeRecord.RequestMeme(req, s, msg)
	case "--ping":
		go c.pingRecord.Pong(s, msg)
	case "--rule34":
		go c.r34Record.Rule34(req, s, msg)
	}
}
