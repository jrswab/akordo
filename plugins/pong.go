package plugins

import (
	"log"

	dg "github.com/bwmarrin/discordgo"
)

// Pong returns the string "Pong" when a user types "--Ping"
func (r *Record) Pong(s *dg.Session, msg *dg.MessageCreate) {
	// Check the last time the user made this request
	if tooSoon := r.checkLastAsk(s, msg); tooSoon {
		return
	}
	_, err := s.ChannelMessageSend(msg.ChannelID, "pong")
	if err != nil {
		log.Fatalf("session.ChannelMessageSend failed: %s", err)
	}
	log.Printf("%s fetched pong", msg.Member.User.Username)
}
