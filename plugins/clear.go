package plugins

import (
	"fmt"
	"log"
	"os"

	dg "github.com/bwmarrin/discordgo"
)

// Eraser is the interface for interacting with the clear package.
type Eraser interface {
	ClearHandler(msg *dg.MessageCreate) error
}

type clear struct {
	dgs *dg.Session
}

// NewEraser creates a new clear struct for using the clear methods
func NewEraser(s *dg.Session) Eraser {
	return &clear{dgs: s}
}

// ClearHandler controls what method is triggered based on the user's command.
func (c *clear) ClearHandler(msg *dg.MessageCreate) error {
	botID, _ := os.LookupEnv("BOT_ID")
	err := c.clearMSGs(botID, msg)
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
	log.Println(messages)

	msgIDs := []string{}
	for _, m := range messages {
		if m.Author.ID == toDeleteID {
			msgIDs = append(msgIDs, m.ID)
		}
	}

	log.Println(msgIDs)
	err = c.dgs.ChannelMessagesBulkDelete(msg.ChannelID, msgIDs)
	if err != nil {

	}
	return nil
}
