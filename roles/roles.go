package roles

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	dg "github.com/bwmarrin/discordgo"
)

// SelfAssignFile is the path were tha self assign role data is located
const SelfAssignFile string = "data/selfAssignRoles.json"

// DgSession is the interface for mocking the discordgo session functions in this package.
type DgSession interface {
	GuildRoles(guildID string) (st []*dg.Role, err error)
	GuildMemberRoleAdd(guildID string, userID string, roleID string) (err error)
	GuildMemberRoleRemove(guildID string, userID string, roleID string) (err error)
}

// Assigner is the interface for interacting with the roleAction methods
type Assigner interface {
	ExecuteRoleCommands(req []string, msg *dg.MessageCreate) (*dg.MessageEmbed, error)
	LoadSelfAssignRoles(file string) error
}
type roleSystem struct {
	dgs DgSession
	sar *roleStorage
}

// Roles holds all data needed to execute the functionality.
type roleStorage struct {
	SelfRoles map[string]string `json:"selfRoles"`
}

// NewRoleStorage creates the roles package data storage map for the session
func NewRoleStorage(s *dg.Session) Assigner {
	return &roleSystem{
		dgs: s,
		sar: &roleStorage{
			SelfRoles: make(map[string]string),
		},
	}
}

func (r *roleSystem) LoadSelfAssignRoles(file string) error {
	savedRoles, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	err = json.Unmarshal(savedRoles, &r.sar.SelfRoles)
	if err != nil {
		return err
	}
	return nil
}

// ExecuteRoles is the method used to run the correct method based on user input.
func (r *roleSystem) ExecuteRoleCommands(req []string, msg *dg.MessageCreate) (*dg.MessageEmbed, error) {
	embed := &dg.MessageEmbed{}
	if len(req) < 2 {
		embed.Description = fmt.Sprintf("Usage: `<prefix>roles lsar`")
		return embed, nil
	}

	switch req[1] {
	case "asar": // add self assign role
		return r.addSelfAssignments(SelfAssignFile, req, msg)
	case "lsar": // list self assignable roles
		return r.listSelfAssignRoles()
	case "sar": // self assign a role
		return r.assignRole(req, msg)
	case "uar":
		return r.unassignRole(req, msg)
	default:
		notListed := "I don't know what to do :thinking:"
		embed.Description = fmt.Sprintf("%s\nPlease check the command and try again", notListed)
	}
	return embed, nil
}

func (r *roleSystem) addSelfAssignments(file string, req []string, msg *dg.MessageCreate) (*dg.MessageEmbed, error) {
	embed := &dg.MessageEmbed{}

	// Look up the bot owner's discord ID
	ownerID, found := os.LookupEnv("BOT_OWNER")
	if !found {
		return embed, fmt.Errorf(
			"roles.go Execute failed: \"BOT_OWNER\" environment variable not found",
		)
	}

	// Make sure the bot owner is running the command
	if msg.Author.ID != ownerID {
		return nil, fmt.Errorf("addSelfAssignments executor is not the bot owner")
	}

	// Tell the user how to use the command if the command is too short
	if len(req) < 3 {
		embed.Description = fmt.Sprintf("Usage: `<prefix>roles asar [role name]`")
		return embed, nil
	}

	// Create a new string if the role has spaces
	sarName := req[2]
	if len(req) > 3 {
		for _, word := range req[3:] {
			sarName = fmt.Sprintf("%s %s", sarName, word)
		}
	}

	// Get roles set in the guild (server)
	roles, err := r.dgs.GuildRoles(msg.GuildID)
	if err != nil {
		return nil, fmt.Errorf("GuildRoles failed: %s", err)
	}

	// Set up the map of role names to their IDs
	for _, role := range roles {
		if role.Name == sarName {
			r.sar.SelfRoles[sarName] = role.ID
		}
	}

	err = r.saveSelfAssignments(file)
	if err != nil {
		return embed, fmt.Errorf("addSelfAssignRole failed: %s", err)
	}

	embed.Description = fmt.Sprintf("Added %s to the self assign role list", sarName)
	return embed, nil
}

func (r *roleSystem) saveSelfAssignments(file string) error {
	json, err := json.MarshalIndent(r.sar.SelfRoles, "", "  ")
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

func (r *roleSystem) listSelfAssignRoles() (*dg.MessageEmbed, error) {
	embed := &dg.MessageEmbed{}
	var assignableRoles string

	for role := range r.sar.SelfRoles {
		assignableRoles = fmt.Sprintf("%s\n%s", assignableRoles, role)
	}

	embed.Description = assignableRoles
	return embed, nil
}

func (r *roleSystem) assignRole(req []string, msg *dg.MessageCreate) (*dg.MessageEmbed, error) {
	embed := &dg.MessageEmbed{}
	sarName, sarID := r.findRequestedRole(req)

	if sarID == "" {
		embed.Description = fmt.Sprintf("Sorry, the role, %s, is not self assignable", sarName)
		return embed, nil
	}

	err := r.dgs.GuildMemberRoleAdd(msg.GuildID, msg.Author.ID, sarID)
	if err != nil {
		return nil, fmt.Errorf("dg method GuildMemberRoleAdd failed: %s", err)
	}

	embed.Description = fmt.Sprintf("Added role, %s, to %s", sarName, msg.Author.Username)
	return embed, nil
}

func (r *roleSystem) unassignRole(req []string, msg *dg.MessageCreate) (*dg.MessageEmbed, error) {
	embed := &dg.MessageEmbed{}
	sarName, sarID := r.findRequestedRole(req)

	if sarID == "" {
		embed.Description = fmt.Sprintf("Sorry, the role, %s, is not self removable", sarName)
		return embed, nil
	}

	err := r.dgs.GuildMemberRoleRemove(msg.GuildID, msg.Author.ID, sarID)
	if err != nil {
		return nil, fmt.Errorf("dg method GuildMemberRoleAdd failed: %s", err)
	}

	embed.Description = fmt.Sprintf("Removed role, %s, from %s", sarName, msg.Author.Username)
	return embed, nil
}

func (r *roleSystem) findRequestedRole(req []string) (string, string) {
	// Create a new string if the role has spaces
	sarName := req[2]
	if len(req) > 3 {
		for _, word := range req[3:] {
			sarName = fmt.Sprintf("%s %s", sarName, word)
		}
	}

	var sarID string
	for role, id := range r.sar.SelfRoles {
		if sarName == role {
			sarID = id
			break
		}
	}
	return sarName, sarID
}
