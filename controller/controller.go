package controller

import (
	"regexp"
	"strings"

	plugs "git.sr.ht/~jrswab/akordo/plugins"
	dg "github.com/bwmarrin/discordgo"
)

// Controller holds the data needed for the bot to send/receive messages.
type Controller struct {
	Session *dg.Session
}

// ReceiveMessage reads the message and returns the response based on the input.
func ReceiveMessage(s *dg.Session, msg *dg.MessageCreate) {
	// Make sure the message matches the bot syntax
	var re = regexp.MustCompile(`(?m)^--(\w|\s)+`)
	match := re.MatchString(msg.Content)
	if !match {
		return
	}

	// Split the string to a slice to parse parameters
	req := strings.Split(msg.Content, " ")

	// Perform action based off each command
	switch req[0] {
	case "--ping":
		go plugs.Pong(s, msg)
	case "--rule34":
		go plugs.Rule34(req, s, msg)
	case "--meme":
		go plugs.RequestMeme(req, s, msg)
	}
}
