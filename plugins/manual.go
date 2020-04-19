package plugins

import (
	"fmt"

	man "git.sr.ht/~jrswab/akordo/manuals"
	dg "github.com/bwmarrin/discordgo"
)

// Manual is triggered when a user passes `--man <command name>`
// This is to give the user a UNIX style man page of what the
// command is able to do.
func Manual(req []string, s *dg.Session, msg *dg.MessageCreate) (string, error) {
	if len(req) < 2 {
		helpMsg := fmt.Sprintf("Usage: `[prefix]man <command name>`")
		return helpMsg, nil
	}

	switch req[1] {
	case "gif":
		return man.Gif, nil
	case "man":
		return man.Man, nil
	case "meme":
		return man.Meme, nil
	case "ping":
		return man.Ping, nil
	case "rule34":
		return man.Rule34, nil
	}
	return "", nil
}
