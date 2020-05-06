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

// AuthClearPath is the path to the json that stores the authorized roles that can  run `clear <username>`
const AuthClearPath string = "data/authorizedToClear.json"

// Eraser is the interface for interacting with the clear package.
type Eraser interface {
	ClearHandler(request []string, msg *dg.MessageCreate) error
	LoadAuthList(file string) error
}

type clear struct {
	dgs   AkSession
	Authd *authorizedRoles `json:"authorizedRoles"`
}

type authorizedRoles struct {
	roleID map[string]bool
}

// NewEraser creates a new clear struct for using the clear methods
func NewEraser(s AkSession) Eraser {
	return &clear{
		dgs: s,
		Authd: &authorizedRoles{
			roleID: make(map[string]bool),
		},
	}
}

// ClearHandler controls what method is triggered based on the user's command.
func (c *clear) ClearHandler(request []string, msg *dg.MessageCreate) error {
	var err error
	var userID string
	if len(request) < 2 {
		botID, _ := os.LookupEnv("BOT_ID")
		err = c.clearMSGs(botID, msg)
	}
	switch request[1] {
	case "set":
		return c.setAuthorized(request[2], msg)
	default:
		userID, err = c.findUserID(request[1], msg)
		err = c.clearMSGs(userID, msg)
	}

	if err != nil {
		return fmt.Errorf("ClearHandler failed: %s", err)
	}

	return nil
}

func (c *clear) setAuthorized(roleID string, msg *dg.MessageCreate) error {
	// Look up the bot owner's discord ID
	ownerID, found := os.LookupEnv(botOwner)
	if !found {
		return fmt.Errorf("clear setAuthorized() failed: %s environment variable not found", botOwner)
	}

	// Make sure the bot owner is executing the command
	if msg.Author.ID != ownerID {
		return nil
	}

	// Make sure the string passed is the role ID
	var re = regexp.MustCompile(atRoleID)
	match := re.MatchString(roleID)
	if !match {
		return fmt.Errorf("Authorized role must be a role and formatted as: `@mod`")
	}

	splitID := strings.Split(roleID, "&")
	id := strings.TrimSuffix(splitID[1], ">")

	c.Authd.roleID[id] = true

	// Save the updates
	err := c.saveAuthList(AuthClearPath)
	if err != nil {
		return fmt.Errorf("error saving role to authorized list: %s", err)
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

func (c *clear) saveAuthList(filePath string) error {
	json, err := json.MarshalIndent(c.Authd, "", "  ")
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

// LoadAuthList loads the saved roles allowed to execute `clear <username>` from the json file
func (c *clear) LoadAuthList(file string) error {
	savedRoles, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	err = json.Unmarshal(savedRoles, &c.Authd.roleID)
	if err != nil {
		return err
	}
	return nil
}
