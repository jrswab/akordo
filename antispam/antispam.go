package antispam

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	dg "github.com/bwmarrin/discordgo"
)

// SpamMax is the path were the bot stores that maximum number of repeated messages before kicking the user
const SpamMax string = "data/spamMax.json"

// SpamTracker is the requierd struct for use with the antispam system.
type SpamTracker struct {
	max      int
	messages map[string][]string
}

// NewSpamTracker creates the requierd struct for use with the antispam system.
func NewSpamTracker() *SpamTracker {
	return &SpamTracker{
		messages: make(map[string][]string),
	}
}

//Handler handles the execution of the antispam methods
func (s *SpamTracker) Handler(request []string, msg *dg.MessageCreate) (bool, error) {
	if len(request) < 1 {
		return false, nil
	}

	if request[1] == "set" {
		// Look up the bot owner's discord ID
		ownerID, found := os.LookupEnv("BOT_OWNER")
		if !found {
			return false, fmt.Errorf("spam set() failed: environment variable not found")
		}

		// Make sure the bot owner is executing the command
		if msg.Author.ID != ownerID {
			return false, nil
		}

		max, err := strconv.Atoi(request[2])
		if err != nil {
			return false, fmt.Errorf("converting string to integer failed: %s", err)
		}
		s.setMax(max)
	}

	return false, nil
}

// CheckForSpam is used to check  if the last `n` messages are the same as the current messag sent by the user.
func (s *SpamTracker) CheckForSpam(msg *dg.MessageCreate) (bool, error) {
	// Add the newest message before checking for spam
	err := s.addMsg(msg)
	if err != nil {
		return false, fmt.Errorf("addMsg() failed: %s", err)
	}

	// Look over the last `n` messages to determine if the user is spamming the same message
	var count int
	for i := (s.max - 1); i > 0; i-- {
		if msg.Content == s.messages[msg.Author.ID][i] {
			count++
		}
		if count == s.max {
			return true, nil
		}
	}
	return false, nil
}

func (s *SpamTracker) addMsg(msg *dg.MessageCreate) error {
	// Check if the user ID is included in the map; if not create it with the current message
	_, found := s.messages[msg.Author.ID]
	if !found {
		s.messages[msg.Author.ID] = []string{msg.Content}
		return nil
	}

	// if the messages to check is at the max, "shift" messages to remove oldest and add newest
	if len(s.messages[msg.Author.ID]) >= s.max {
		newMsgList := []string{}
		for idx, val := range s.messages[msg.Author.ID] {
			if idx == 0 {
				continue
			}
			newMsgList = append(newMsgList, val)
		}
	}

	// Append the most recent message
	s.messages[msg.Author.ID] = append(s.messages[msg.Author.ID], msg.Content)

	return nil
}

// setMax sets the maximum number of repeated messages before the user is kicked from the chat.
func (s *SpamTracker) setMax(n int) error {
	s.max = n
	err := s.saveMax(SpamMax)
	if err != nil {
		return fmt.Errorf("setMax could not save the requested maximum: %s", err)
	}
	return nil
}

func (s *SpamTracker) saveMax(file string) error {
	json, err := json.MarshalIndent(s.max, "", "  ")
	if err != nil {
		return err
	}

	// Write to data to a file
	err = ioutil.WriteFile(file, json, 0600)
	if err != nil {
		return err
	}
	return nil
}

// LoadMax loads the saved maximum value set by the bot owner
func (s *SpamTracker) LoadMax(file string) error {
	savedMax, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	err = json.Unmarshal(savedMax, &s.max)
	if err != nil {
		return err
	}
	return nil
}
