package plugins

import (
	"fmt"

	man "git.sr.ht/~jrswab/akordo/manuals"
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
		return man.Gif
	case "man":
		return man.Man
	case "meme":
		return man.Meme
	case "ping":
		return man.Ping
	case "rule34":
		return man.Rule34
	}
	return ""
}
