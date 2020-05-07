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

// Agreement contains the data needed to save and execute the base role
type Agreement struct {
	session  AkSession
	BaseRole string `json:"baseRole"` // ID of the role to give
}

// NewAgreement creates the struct to use the methods in rules.go
func NewAgreement(s *dg.Session) *Agreement {
	return &Agreement{
		session:  s,
		BaseRole: "",
	}
}

// Handler triggers the correct method based off the user input.
func (a *Agreement) Handler(req []string, msg *dg.MessageCreate) (string, error) {
	if len(req) < 2 {
		return "Usage: <prefix>rules agreed", nil
	}

	switch req[1] {
	case "set":
		return a.addAgreementRole(req[2], msg)
	case "agreed":
		return a.ruleAgreement(msg)
	}

	return "Usage: <prefix>rules agreed", nil
}

func (a *Agreement) ruleAgreement(msg *dg.MessageCreate) (string, error) {
	err := a.session.GuildMemberRoleAdd(msg.GuildID, msg.Author.ID, a.BaseRole)
	if err != nil {
		return "", fmt.Errorf("GuildMemberRoleAdd() failed: %s", err)
	}

	return "Added :ok_hand:", nil
}

func (a *Agreement) addAgreementRole(roleID string, msg *dg.MessageCreate) (string, error) {
	// Look up the bot owner's discord ID
	ownerID, found := os.LookupEnv("BOT_OWNER")
	if !found {
		return "", fmt.Errorf(
			"banned words addEditors() failed: \"BOT_OWNER\" environment variable not found",
		)
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

	a.BaseRole = id

	err := a.saveRole(ChatPermissionRole)
	if err != nil {
		return "", fmt.Errorf("error saving base role: %s", err)
	}
	return "Role added as as the base chat role", nil
}

func (a *Agreement) saveRole(filePath string) error {
	json, err := json.MarshalIndent(a.BaseRole, "", "  ")
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

// LoadAgreementRole loads the saved json into the struct of banned words
func (a *Agreement) LoadAgreementRole(filePath string) error {
	baseRole, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	err = json.Unmarshal(baseRole, &a.BaseRole)
	if err != nil {
		return err
	}
	return nil
}
