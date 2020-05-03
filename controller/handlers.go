package controller

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	plugs "gitlab.com/technonauts/akordo/plugins"
)

func (c *controller) checkWords() {
	sd := c.sess
	// Check for blacklisted words
	isBlacklisted, err := sd.Blacklist.CheckBannedWords(c.msg)
	if err != nil {
		log.Printf("CheckBannedWords() failed: %s", err)
	}
	if isBlacklisted {
		reason := "Kicked for inappropriate language."
		err := sd.session.GuildMemberDeleteWithReason(c.msg.GuildID, c.msg.Author.ID, reason)
		if err != nil {
			log.Printf("GuildMemberDeleteWithReason() failed: %s", err)
		}
	}
}

func (c *controller) determineIfCmd() bool {
	sd := c.sess
	// Make sure the message matches the bot syntax
	regEx := fmt.Sprintf("(?m)^%s(\\w|\\s)+", sd.prefix)
	var re = regexp.MustCompile(regEx)
	match := re.MatchString(c.msg.Content)
	if match {
		return true
	}

	exempt := xpExemptions(c.msg.Content)
	if exempt {
		return false
	}

	// Add xp for all non-bot messages
	sd.XP.ManipulateXP("addMessagePoints", c.msg)

	// Check for role promotion
	err := sd.Roles.AutoPromote(c.msg)
	if err != nil {
		log.Printf("xp.AutoPromote failed: %s", err)
	}

	return false
}

// Handler looks up the command found by the bot and kicks off a Goroutine do what
// the user is asking to do.
//
// To remove a plugin simply remove the case statement for that plugin
// To add a plugin, create a case statement for the plugin as shown below.
// If the plugin is new create a new `.go` file under the `plugins` directory.
func (c *controller) cmdHandler() {
	sd := c.sess
	msg := c.msg
	var err error

	// Split the string to a slice to parse parameters
	req := strings.Split(msg.Content, " ")

	switch req[0] {
	case sd.prefix + "blacklist":
		c.response, err = sd.Blacklist.Handler(req, msg)

	case sd.prefix + "clear":
		c.msgType = "none"
		err = sd.clear.ClearHandler(msg)

	case sd.prefix + "crypto":
		c.response, err = sd.crypto.Game(req, msg)

	case sd.prefix + "gif":
		c.response, err = sd.gifRequest.Gif(req, sd.session, msg)

	case sd.prefix + "man":
		c.msgType = "dm"
		c.response = plugs.Manual(req, sd.session, msg)

	case sd.prefix + "meme":
		c.response, err = sd.memeRequest.RequestMeme(req, sd.session, msg)

	case sd.prefix + "ping":
		c.response = sd.pingRecord.Pong(msg)

	case sd.prefix + "roles":
		c.msgType = "embed"
		c.emb, err = sd.Roles.ExecuteRoleCommands(req, msg)

	case sd.prefix + "rule34":
		c.response, err = sd.r34Request.Rule34(req, sd.session, msg)

	case sd.prefix + "rules":
		c.delete = true
		c.response, err = sd.Rules.Handler(req, msg)

	case sd.prefix + "version":
		c.msgType = "embed"
		c.emb, err = printVersion()

	case sd.prefix + "xp":
		c.msgType = "embed"
		c.emb, err = sd.XP.Execute(req, msg)

	default:
		c.response = "I don't know what to do with that :sob:"
	}

	if err != nil {
		log.Printf("error executing task: %s", err)
		return
	}

	c.reply()
}
