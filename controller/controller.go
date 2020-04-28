package controller

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"sync"

	plugs "git.sr.ht/~jrswab/akordo/plugins"
	"git.sr.ht/~jrswab/akordo/roles"
	"git.sr.ht/~jrswab/akordo/xp"
	dg "github.com/bwmarrin/discordgo"
)

// SessionData holds the data needed to complete the requested transactions
type SessionData struct {
	session *dg.Session
	Mutex   *sync.Mutex
	prefix  string
	XP      xp.Exp
	Roles   roles.Assigner

	// Plugins:
	Blacklist   *plugs.Blacklist
	clear       plugs.Eraser
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
		session: s,
		Mutex:   &sync.Mutex{},
		prefix:  `=`,

		// Plugins:
		crypto:      plugs.NewCrypto(),
		gifRequest:  plugs.NewGifRequest(),
		memeRequest: plugs.NewMemeRequest(),
		pingRecord:  plugs.NewRecorder(),
		r34Request:  plugs.NewRule34Request(),
	}

	// Commands that require the session for execution
	sd.XP = xp.NewXpStore(sd.Mutex, sd.session)
	sd.Roles = roles.NewRoleStorage(sd.session)
	sd.clear = plugs.NewEraser(sd.session)
	sd.Blacklist = plugs.NewBlacklist(sd.session)
	sd.Rules = plugs.NewAgreement(sd.session)

	return sd
}

// NewMessage waits for a ne message to be sent in a the Discord guild
// This kicks off a Goroutine to free up the mutex set by discordgo `AddHandler` method.
func (sd *SessionData) NewMessage(s *dg.Session, msg *dg.MessageCreate) {
	go sd.checkMessage(msg)
}

// CheckSyntax uses regexp from the standard library to check the message has the correct
// prefix as defined by the `prefix` constant.
func (sd *SessionData) checkMessage(msg *dg.MessageCreate) {
	// Check for blacklisted words
	isBlacklisted, err := sd.Blacklist.CheckBannedWords(msg)
	if err != nil {
		log.Printf("CheckBannedWords() failed: %s", err)
	}
	if isBlacklisted {
		reason := "Kicked for inappropriate language."
		err := sd.session.GuildMemberDeleteWithReason(msg.GuildID, msg.Author.ID, reason)
		if err != nil {
			log.Printf("GuildMemberDeleteWithReason() failed: %s", err)
		}
	}

	// Make sure the message matches the bot syntax
	regEx := fmt.Sprintf("(?m)^%s(\\w|\\s)+", sd.prefix)
	var re = regexp.MustCompile(regEx)
	match := re.MatchString(msg.Content)
	if !match {
		// Add xp for all non-bot messages
		sd.XP.ManipulateXP("addMessagePoints", msg)
		// Check for role promotion
		err := sd.XP.AutoPromote(msg)
		if err != nil {
			log.Printf("xp.AutoPromote failed: %s", err)
		}
		return
	}

	sd.ExecuteTask(msg)

	// Remove bot command after the bot replies
	err = sd.session.ChannelMessageDelete(msg.ChannelID, msg.ID)
	if err != nil {
		log.Printf("failed to delete message after bot reply: %s", err)
	}
}

// ExecuteTask looks up the command found by the bot and kicks off a Goroutine do what
// the user is asking to do.
//
// To remove a plugin simply remove the case statement for that plugin
// To add a plugin, create a case statement for the plugin as shown below.
// If the plugin is new create a new `.go` file under the `plugins` directory.
func (sd *SessionData) ExecuteTask(msg *dg.MessageCreate) {
	var (
		res string
		emb *dg.MessageEmbed
		err error
	)

	// Split the string to a slice to parse parameters
	req := strings.Split(msg.Content, " ")

	msgType := "chan"
	switch req[0] {
	case sd.prefix + "blacklist":
		res, err = sd.Blacklist.Handler(req, msg)
	case sd.prefix + "clear":
		msgType = "none"
		err = sd.clear.ClearHandler(msg)
	case sd.prefix + "crypto":
		res, err = sd.crypto.Game(req, msg)
	case sd.prefix + "gif":
		res, err = sd.gifRequest.Gif(req, sd.session, msg)
	case sd.prefix + "man":
		msgType = "dm"
		res = plugs.Manual(req, sd.session, msg)
	case sd.prefix + "meme":
		res, err = sd.memeRequest.RequestMeme(req, sd.session, msg)
	case sd.prefix + "ping":
		res = sd.pingRecord.Pong(msg)
	case sd.prefix + "roles":
		msgType = "embed"
		emb, err = sd.Roles.ExecuteRoleCommands(req, msg)
	case sd.prefix + "rule34":
		res, err = sd.r34Request.Rule34(req, sd.session, msg)
	case sd.prefix + "rules":
		res, err = sd.Rules.Handler(req, msg)
	case sd.prefix + "xp":
		res, err = sd.XP.Execute(req, msg)
	default:
		res = "I don't know what to do with that :sob:"
	}

	if err != nil {
		log.Printf("error executing task: %s", err)
		return
	}

	sd.reply(res, msgType, emb, msg)
}

// Reply takes the executed data and replies to the user. This is either in the channel
// where the command was sent or as a direct message to the user.
func (sd *SessionData) reply(res, msgType string, emb *dg.MessageEmbed, msg *dg.MessageCreate) {
	s := sd.session

	switch msgType {
	case "dm":
		sd.sendAsDM(res, msg)
	case "embed":
		_, err := s.ChannelMessageSendEmbed(msg.ChannelID, emb)
		if err != nil {
			log.Printf("reply ChannelMessageSendEmbed failed: %s", err)
		}
	case "chan":
		_, err := s.ChannelMessageSend(msg.ChannelID, res)
		if err != nil {
			log.Printf("reply ChannelMessageSend failed: %s", err)
		}
	}
}

func (sd *SessionData) sendAsDM(res string, msg *dg.MessageCreate) {
	s := sd.session
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
