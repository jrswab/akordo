package plugins

import (
	"fmt"
	"time"

	dg "github.com/bwmarrin/discordgo"
)

// CommandDelay is the time (in seconds) to restrict command spam. Exported for unit tests.
const CommandDelay = 90
const botDelay = CommandDelay * time.Second

// ChatPermissionRole is the path to the banned words json file
const ChatPermissionRole string = "data/ruleRole.json"

// AuthClearPath is the path to the json that stores the authorized roles that can  run `clear <username>`
const AuthClearPath string = "data/authorizedToClear.json"

// BannedWordsPath is the path to the banned words json file
const BannedWordsPath string = "data/bannedWords.json"

// Environment variables for the plugins package.
const botOwner = "BOT_OWNER"

const atRoleID string = `(?m)^<@&\d+>$`

// AkSession allows for tests to mock the discordgo session.Channel() method call
type AkSession interface {
	Channel(channelID string) (st *dg.Channel, err error)
	GuildMemberRoleAdd(guildID, userID, roleID string) (err error)
	ChannelMessages(channelID string, limit int, beforeID, afterID, aroundID string) (st []*dg.Message, err error)
	ChannelMessagesBulkDelete(channelID string, messages []string) (err error)
	GuildMembers(guildID string, after string, limit int) (st []*dg.Member, err error)
}

// Record holds the users last gif request to avoid spamming.
type Record struct {
	MinWaitTime time.Duration
	LastReq     map[string]time.Time
}

// NewRecorder creates a Recordger with a defined map
func NewRecorder() *Record {
	userMap := make(map[string]time.Time)
	return &Record{LastReq: userMap, MinWaitTime: (botDelay)}
}

// CheckLastAsk checks the last time the user executed the specific command
func (r *Record) CheckLastAsk(msg *dg.MessageCreate) (string, bool) {
	last, found := r.LastReq[msg.Author.ID]
	if found && time.Since(last) < (r.MinWaitTime) {
		userAlert := fmt.Sprintf("%s please wait %d seconds before requesting the same command.",
			msg.Author.Username, CommandDelay)

		return userAlert, true
	}

	// Add or update the new request
	r.LastReq[msg.Author.ID] = time.Now()
	return "", false
}
