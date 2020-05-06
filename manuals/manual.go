package manuals

import (
	"fmt"

	dg "github.com/bwmarrin/discordgo"
)

// Manual is triggered when a user passes `prefix man <command name>`
// This is to give the user a UNIX style man page of what the
// command is able to do.
func Manual(req []string, s *dg.Session, msg *dg.MessageCreate) string {
	if len(req) < 2 {
		helpMsg := fmt.Sprintf("Usage: `<prefix>man command`")
		return helpMsg
	}

	switch req[1] {
	case "gif":
		return gif
	case "man":
		return man
	case "meme":
		return meme
	case "ping":
		return ping
	case "rule34":
		return rule34
	}
	return "Sorry, I don't have a manual for that :confused:"
}
