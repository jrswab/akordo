package plugins

import (
	"fmt"
	"log"
	"time"

	dg "github.com/bwmarrin/discordgo"
)

// Record holds the users last gif request to avoid spamming.
type Record struct {
	minWaitTime time.Duration
	lastReq     map[string]time.Time
}

// NewRecorder creates a Recordger with a defined map
func NewRecorder() *Record {
	userMap := make(map[string]time.Time)
	return &Record{lastReq: userMap, minWaitTime: (2 * time.Minute)}
}

func (r *Record) checkLastAsk(s *dg.Session, msg *dg.MessageCreate) bool {
	last, found := r.lastReq[msg.Author.ID]
	if found && time.Since(last) < (r.minWaitTime) {
		userMention := fmt.Sprintf("%s please wait 120 seconds before requesting another Gif.",
			msg.Author.Username)

		_, err := s.ChannelMessageSend(msg.ChannelID, userMention)
		if err != nil {
			log.Printf("session.ChannelMessageSend failed: %s", err)
		}
		return true
	}

	// Add or update the new request
	r.lastReq[msg.Author.ID] = time.Now()
	return false
}
