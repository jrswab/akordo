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
func (s *SpamTracker) Handler(request []string, msg *dg.MessageCreate) (string, error) {
	if len(request) < 2 {
		return "Usage: <prefix>antispam set [number]", nil
	}

	if request[1] != "set" {
		return "Usage: <prefix>antispam set [number]", nil
	}
	// Look up the bot owner's discord ID
	ownerID, found := os.LookupEnv("BOT_OWNER")
	if !found {
		return "", fmt.Errorf("spam set() failed: environment variable not found")
	}

	// Make sure the bot owner is executing the command
	if msg.Author.ID != ownerID {
		return "", nil
	}

	max, err := strconv.Atoi(request[2])
	if err != nil {
		return "", fmt.Errorf("converting string to integer failed: %s", err)
	}
	s.setMax(max)

	return "Max repeated messages for kick set.", nil
}

// CheckForSpam is used to check  if the last `n` messages are the same as the current messag sent by the user.
func (s *SpamTracker) CheckForSpam(msg *dg.MessageCreate) (bool, error) {
	// If max is zero don't checking
	if s.max == 0 {
		return false, nil
	}

	// Add the newest message before checking for spam
	isNewMsg := s.addMsg(msg)
	if isNewMsg {
		return false, nil
	}

	// Look over the last `n` messages to determine if the user is spamming the same message
	var count int
	for i := (s.max - 1); i >= 0; i-- {
		if msg.Content == s.messages[msg.Author.ID][i] {
			count++
		}
		if count == s.max {
			return true, nil
		}
	}
	return false, nil
}

func (s *SpamTracker) addMsg(msg *dg.MessageCreate) bool {
	// Check if the user ID is included in the map; if not create it with the current message
	_, found := s.messages[msg.Author.ID]
	if !found {
		s.messages[msg.Author.ID] = []string{msg.Content}
		return true
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

	return false
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
