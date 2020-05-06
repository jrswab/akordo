package plugins

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	dg "github.com/bwmarrin/discordgo"
)

// BannedWordsPath is the path to the banned words json file
const BannedWordsPath string = "data/bannedWords.json"

// Blacklist contains the data needed to execute the method functionality
type Blacklist struct {
	session AkSession
	data    *data
}

type data struct {
	// Role IDs that are allowed to add/remove banned words
	Editors map[string]bool `json:"editors"`
	Banned  map[string]bool `json:"blacklist"`
}

// NewBlacklist creates the needed structs for running the ban methods
func NewBlacklist(s AkSession) *Blacklist {
	return &Blacklist{
		session: s,
		data: &data{
			Editors: make(map[string]bool),
			Banned:  make(map[string]bool),
		},
	}
}

// Handler executes the correct method based off the users message
func (b *Blacklist) Handler(req []string, msg *dg.MessageCreate) (string, error) {
	switch req[1] {
	case "add":
		return b.addBlacklistWord(req, msg)
	case "remove":
		return b.removeBlacklistWord(req, msg)
	case "editor":
		return b.addedEditors(req[2], msg)
	}
	return "I don't know what to do with that :sob:", nil
}

// CheckBannedWords checks a message for banned words if the channel is not marked NSFW
func (b *Blacklist) CheckBannedWords(msg *dg.MessageCreate) (bool, error) {
	// Ignore banned words sent by an editor (so they don't get banned if removing one)
	for _, role := range msg.Member.Roles {
		if b.data.Editors[role] {
			return false, nil
		}
	}

	// Ignore bot owner messages
	if msg.Author.ID == os.Getenv(botOwner) {
		return false, nil
	}

	// Ignore if the channel is marked NSFW
	dChan, err := b.session.Channel(msg.ChannelID)
	if err != nil {
		return false, fmt.Errorf("discordGo Channel() failed: %s", err)
	}

	if !dChan.NSFW {
		words := strings.Split(msg.Content, " ")
		for _, word := range words {
			if b.data.Banned[word] {
				return true, nil
			}
		}
	}
	return false, nil
}

func (b *Blacklist) addBlacklistWord(req []string, msg *dg.MessageCreate) (string, error) {
	word := strings.Join(req[2:], " ")
	b.data.Banned[word] = true

	err := b.saveWordList(BannedWordsPath)
	if err != nil {
		return "", fmt.Errorf("error saving new word to list: %s", err)
	}

	return "Word(s) added to the blacklist", nil
}

func (b *Blacklist) removeBlacklistWord(req []string, msg *dg.MessageCreate) (string, error) {
	for _, word := range req[2:] {
		b.data.Banned[word] = false
	}

	err := b.saveWordList(BannedWordsPath)
	if err != nil {
		return "", fmt.Errorf("error saving word removal to list: %s", err)
	}

	return "Word(s) removed from the blacklist", nil
}

func (b *Blacklist) saveWordList(filePath string) error {
	json, err := json.MarshalIndent(b.data, "", "  ")
	if err != nil {
		return err
	}

	// Write to data to a file
	err = ioutil.WriteFile(filePath, json, 0600)
	if err != nil {
		return err
	}
	return nil
}

// LoadBannedWordList loads the saved json into the struct of banned words
func (b *Blacklist) LoadBannedWordList(filePath string) error {
	savedBannedWords, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	err = json.Unmarshal(savedBannedWords, &b.data)
	if err != nil {
		return err
	}
	return nil
}

func (b *Blacklist) addedEditors(roleID string, msg *dg.MessageCreate) (string, error) {
	// Look up the bot owner's discord ID
	ownerID, found := os.LookupEnv(botOwner)
	if !found {
		return "", fmt.Errorf("banned words addEditors() failed: %s environment variable not found", botOwner)
	}

	// Make sure the bot owner is executing the command
	if msg.Author.ID != ownerID {
		return "This command is for the bot owner only :rage:", nil
	}

	// Make sure the string passed is the role ID
	var re = regexp.MustCompile(`(?m)^<@&\d+>$`)
	match := re.MatchString(roleID)
	if !match {
		return "Editor must be a role and formatted as: `@mod`", nil
	}

	splitID := strings.Split(roleID, "&")
	id := strings.TrimSuffix(splitID[1], ">")

	b.data.Editors[id] = true

	// Save the updates
	err := b.saveWordList(BannedWordsPath)
	if err != nil {
		return "", fmt.Errorf("error saving editor to word list: %s", err)
	}
	return "Role added as an editor", nil
}
