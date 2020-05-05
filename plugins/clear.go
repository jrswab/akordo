package plugins

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	dg "github.com/bwmarrin/discordgo"
)

// Eraser is the interface for interacting with the clear package.
type Eraser interface {
	ClearHandler(request []string, msg *dg.MessageCreate) error
}

type clear struct {
	dgs *dg.Session
}

// NewEraser creates a new clear struct for using the clear methods
func NewEraser(s *dg.Session) Eraser {
	return &clear{dgs: s}
}

// ClearHandler controls what method is triggered based on the user's command.
func (c *clear) ClearHandler(request []string, msg *dg.MessageCreate) error {
	var err error
	var userID string
	switch len(request) {
	case 1:
		botID, _ := os.LookupEnv("BOT_ID")
		err = c.clearMSGs(botID, msg)
	case 2:
		userID, err = c.findUserID(request[1], msg)
		err = c.clearMSGs(userID, msg)
	}

	if err != nil {
		return fmt.Errorf("ClearHandler failed: %s", err)
	}

	return nil
}

func (c *clear) clearMSGs(toDeleteID string, msg *dg.MessageCreate) error {
	messages, err := c.dgs.ChannelMessages(msg.ChannelID, 100, "", "", "")
	if err != nil {
		return fmt.Errorf("clearMSGs failed: %s", err)
	}

	msgIDs := []string{}
	for _, m := range messages {
		if m.Author.ID == toDeleteID {
			msgIDs = append(msgIDs, m.ID)
		}
	}

	err = c.dgs.ChannelMessagesBulkDelete(msg.ChannelID, msgIDs)
	if err != nil {
		return fmt.Errorf("ChannelMessageBulkDelete failed: %s", err)
	}
	return nil
}

func (c *clear) findUserID(userName string, msg *dg.MessageCreate) (string, error) {
	var (
		id      string
		members []*dg.Member
		err     error
	)

	var re = regexp.MustCompile("(?m)^<@!\\w+>")
	match := re.MatchString(userName)
	if match {
		id = strings.TrimPrefix(userName, "<@!")
		id = strings.TrimSuffix(id, ">")
		return id, nil
	}

	// Get first round of members
	current, err := c.dgs.GuildMembers(msg.GuildID, "", 1000)
	if err != nil {
		return "", err
	}

	for _, names := range current {
		members = append(members, names)
	}

	// if first round has 1000 entries run again until all members are present.
	for len(current) == 1000 {
		lastMember := current[len(current)-1]
		current, err = c.dgs.GuildMembers(msg.GuildID, lastMember.User.ID, 1000)
		if err != nil {
			return "", fmt.Errorf("GuildMembers failed: %s", err)
		}

		for _, names := range current {
			members = append(members, names)
		}

	}

	for _, member := range members {
		if member.Nick == userName {
			id = member.User.ID
			break
		}

		if member.User.Username == userName {
			id = member.User.ID
			break
		}
	}

	return id, nil
}
