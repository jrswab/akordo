package controller

import (
	"log"
	"regexp"
	"sync"

	dg "github.com/bwmarrin/discordgo"
	"gitlab.com/technonauts/akordo/antispam"
	plugs "gitlab.com/technonauts/akordo/plugins"
	"gitlab.com/technonauts/akordo/roles"
	"gitlab.com/technonauts/akordo/xp"
)

const version string = "v1.0.0"

// SessionData holds the data needed to complete the requested transactions
type SessionData struct {
	session  *dg.Session
	Mutex    *sync.Mutex
	prefix   string
	XP       *xp.System
	Roles    roles.Assigner
	AntiSpam *antispam.SpamTracker

	// Plugins:
	Blacklist   *plugs.Blacklist
	Clear       plugs.Eraser
	crypto      *plugs.Crypto
	gifRequest  *plugs.GifRequest
	memeRequest *plugs.MemeRequest
	pingRecord  *plugs.Record
	r34Request  *plugs.Rule34Request
	Rules       *plugs.Agreement
}

// NewSessionData creates a SessionData
func NewSessionData(s *dg.Session) *SessionData {
	sd := &SessionData{
		session:  s,
		Mutex:    &sync.Mutex{},
		prefix:   `=`,
		AntiSpam: antispam.NewSpamTracker(),

		// Plugins:
		crypto:      plugs.NewCrypto(),
		gifRequest:  plugs.NewGifRequest(),
		memeRequest: plugs.NewMemeRequest(),
		pingRecord:  plugs.NewRecorder(),
		r34Request:  plugs.NewRule34Request(),
	}

	// Commands that require the session for execution
	sd.XP = xp.NewXpStore(sd.Mutex, sd.session)
	sd.Roles = roles.NewRoleStorage(sd.session, sd.XP)
	sd.Clear = plugs.NewEraser(sd.session)
	sd.Blacklist = plugs.NewBlacklist(sd.session)
	sd.Rules = plugs.NewAgreement(sd.session)

	return sd
}

type controller struct {
	sess     *SessionData
	msgType  string
	response string
	delete   bool
	emb      *dg.MessageEmbed
	msg      *dg.MessageCreate
}

// NewMessage waits for a ne message to be sent in a the Discord guild
// This kicks off a Goroutine to free up the mutex set by discordgo `AddHandler` method.
func (sd *SessionData) NewMessage(s *dg.Session, msg *dg.MessageCreate) {
	// Create a new controller for each message containing the data used acrossed
	// the entire session.
	c := &controller{
		sess:     sd,
		msgType:  "chan",
		response: "",
		emb:      &dg.MessageEmbed{},
		msg:      msg,
	}

	go c.checkMessage()
}

// CheckSyntax uses regexp from the standard library to check the message has the correct
// prefix as defined by the `prefix` constant.
func (c *controller) checkMessage() {
	// If the message is not a command, award the user XP and ignore the rest of the function
	isCMD := c.determineIfCmd()
	if isCMD {
		// Execute user command:
		c.cmdHandler()

		// Remove command after the bot replies
		c.deleteMessage()
		return
	}

	// Check for spamming
	isSpam, err := c.sess.AntiSpam.CheckForSpam(c.msg)
	if err != nil {
		log.Printf("CheckForSpam failed: %s\n", err)
	}

	if isSpam {
		// Kick the offender
		reason := "Kicked for spamming the chat."
		err := c.sess.session.GuildMemberDeleteWithReason(c.msg.GuildID, c.msg.Author.ID, reason)
		if err != nil {
			log.Printf("isSpam: GuildMemberDeleteWithReason() failed: %s\n", err)
		}
		return
	}

	// Check for banned words in the message.
	foundBannedWord, err := c.checkWords()
	if err != nil {
		log.Printf("checkWords failed: %s\n", err)
	}

	// Stop execution if a banned word is found
	if foundBannedWord {
		// Remove the message
		c.deleteMessage()
		return
	}

	// Award xp if all checks pass
	c.awardXP()
	return
}

func (c *controller) deleteMessage() error {
	sd := c.sess

	err := sd.session.ChannelMessageDelete(c.msg.ChannelID, c.msg.ID)
	if err != nil {
		log.Printf("failed to delete message after bot reply: %s\n", err)
	}
	return nil
}

// Reply takes the executed data and replies to the user. This is either in the channel
// where the command was sent or as a direct message to the user.
func (c *controller) reply() {
	sd := c.sess
	s := sd.session
	var (
		botMsg *dg.Message
		err    error
	)

	// Determine bot output method.
	switch c.msgType {
	case "dm":
		c.sendAsDM()
	case "embed":
		botMsg, err = s.ChannelMessageSendEmbed(c.msg.ChannelID, c.emb)
		if err != nil {
			log.Printf("reply ChannelMessageSendEmbed failed: %s\n", err)
		}
	case "chan":
		botMsg, err = s.ChannelMessageSend(c.msg.ChannelID, c.response)
		if err != nil {
			log.Printf("reply ChannelMessageSend failed: %s\n", err)
		}
	}

	// Delete bot reply if true
	if c.delete {
		err = sd.session.ChannelMessageDelete(botMsg.ChannelID, botMsg.ID)
		if err != nil {
			log.Printf("failed to delete message after bot reply: %s\n", err)
		}
	}
}

func (c *controller) sendAsDM() {
	sd := c.sess
	s := sd.session
	dm, err := s.UserChannelCreate(c.msg.Author.ID)
	if err != nil {
		log.Printf("s.UserChannelCreate failed to create DM for %s: %s\n",
			c.msg.Author.Username, err)
		return
	}

	_, err = s.ChannelMessageSend(dm.ID, c.response)
	if err != nil {
		log.Printf("session.ChannelMessageSend failed to send DM: %s\n", err)
	}
	return
}

func xpExemptions(msg string) bool {
	// Check for emoji only message
	var re = regexp.MustCompile(`(?m)^:(.)+:$`)
	match := re.MatchString(msg)
	if match {
		return true
	}

	// Check for ten or more repeating characters
	var prevByte rune
	var count uint
	for _, v := range msg {

		if v == prevByte {
			count++
		}

		if count >= 10 {
			return true
		}

		prevByte = v
	}
	return false
}

func printVersion() (*dg.MessageEmbed, error) {
	return &dg.MessageEmbed{Title: "Current Version:", Description: version}, nil
}
