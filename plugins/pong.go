package plugins

import (
	"log"

	dg "github.com/bwmarrin/discordgo"
)

// Pong returns the string "Pong" when a user types "--Ping"
func Pong(s *dg.Session, msg *dg.MessageCreate) {
	_, err := s.ChannelMessageSend(msg.ChannelID, "pong")
	if err != nil {
		log.Fatalf("session.ChannelMessageSend failed: %s", err)
	}
	log.Printf("%v fetched pong", msg.Member)
}
