package controller

import (
	"log"

	dg "github.com/bwmarrin/discordgo"
)

// Controller holds the data needed for the bot to send/receive messages.
type Controller struct {
	Session *dg.Session
}

// ReceiveMessage reads the message and returns the response based on the input.
func ReceiveMessage(s *dg.Session, msg *dg.MessageCreate) {
	if msg.Content == "--ping" {
		_, err := s.ChannelMessageSend(msg.ChannelID, "pong")
		if err != nil {
			log.Fatalf("session.ChannelMessageSend failed: %s", err)
		}
	}
}
