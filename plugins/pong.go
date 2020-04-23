package plugins

import (
	dg "github.com/bwmarrin/discordgo"
)

// Pong returns the string "Pong" when a user types "--Ping"
func (r *Record) Pong(msg *dg.MessageCreate) string {
	// Check the last time the user made this request
	alertUser, tooSoon := r.CheckLastAsk(msg)
	if tooSoon {
		return alertUser
	}

	return "pong"
}
