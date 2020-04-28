package plugins

import (
	"fmt"
	"time"

	dg "github.com/bwmarrin/discordgo"
)

// AkSession allows for tests to mock the discordgo session.Channel() method call
type AkSession interface {
	Channel(channelID string) (st *dg.Channel, err error)
	GuildMemberRoleAdd(guildID, userID, roleID string) (err error)
}

// Record holds the users last gif request to avoid spamming.
type Record struct {
	MinWaitTime time.Duration
	LastReq     map[string]time.Time
}

// NewRecorder creates a Recordger with a defined map
func NewRecorder() *Record {
	userMap := make(map[string]time.Time)
	return &Record{LastReq: userMap, MinWaitTime: (2 * time.Minute)}
}

func (r *Record) CheckLastAsk(msg *dg.MessageCreate) (string, bool) {
	last, found := r.LastReq[msg.Author.ID]
	if found && time.Since(last) < (r.MinWaitTime) {
		userAlert := fmt.Sprintf("%s please wait 120 seconds before requesting the same command.",
			msg.Author.Username)

		return userAlert, true
	}

	// Add or update the new request
	r.LastReq[msg.Author.ID] = time.Now()
	return "", false
}
