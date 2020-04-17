package plugins

import (
	"log"

	man "git.sr.ht/~jrswab/akordo/manuals"
	dg "github.com/bwmarrin/discordgo"
)

// Manual is triggered when a user passes `--man <command name>`
// This is to give the user a UNIX style man page of what the
// command is able to do.
func Manual(req []string, s *dg.Session, msg *dg.MessageCreate) {
	memberID := msg.Member.User.ID
	if len(req) < 2 {
		_, err := s.ChannelMessageSend(msg.ChannelID, "Usage: `--man <command name>`")
		if err != nil {
			log.Printf("session.ChannelMessageSend failed: %s", err)
		}
		return
	}

	switch req[1] {
	case "rule34":
		dm, err := s.UserChannelCreate(memberID)
		if err != nil {
			log.Printf("s.UserChannelCreate failed to create DM for %s: %s", memberID, err)
			return
		}
		_, err = s.ChannelMessageSend(dm.ID, man.Rule34)
		if err != nil {
			log.Printf("session.ChannelMessageSend failed: %s", err)
		}
	}
}
