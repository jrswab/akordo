package roles

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strconv"

	"git.sr.ht/~jrswab/akordo/xp"
	x "git.sr.ht/~jrswab/akordo/xp"
	dg "github.com/bwmarrin/discordgo"
)

// SelfAssignFile is the path were tha self assign role data is located
const SelfAssignFile string = "data/selfAssignRoles.json"

// MsgEmbed is used to shorten the name of the original embed type from discordGo
type MsgEmbed *dg.MessageEmbed

// DgSession is the interface for mocking the discordgo session functions in this package.
type DgSession interface {
	GuildRoles(guildID string) (st []*dg.Role, err error)
	GuildMember(guildID string, userID string) (st *dg.Member, err error)
	GuildMembers(guildID string, after string, limit int) (st []*dg.Member, err error)
	GuildMemberRoleAdd(guildID string, userID string, roleID string) (err error)
	GuildMemberRoleRemove(guildID string, userID string, roleID string) (err error)
}

// Assigner is the interface for interacting with the roleAction methods
type Assigner interface {
	ExecuteRoleCommands(req []string, msg *dg.MessageCreate) (*dg.MessageEmbed, error)
	LoadSelfAssignRoles(file string) error
	AutoPromote(msg *dg.MessageCreate) error
	LoadAutoRanks(file string) error
}
type roleSystem struct {
	dgs   DgSession
	xp    *xp.System
	sar   *roleStorage
	tiers *autoRanks
}

// Roles holds all data needed to execute the functionality.
type roleStorage struct {
	SelfRoles map[string]string `json:"selfRoles"`
}

type autoRanks struct {
	Tiers map[string]float64 `json:"tiers"` // map of total xp of role IDs
}

// NewRoleStorage creates the roles package data storage map for the session
func NewRoleStorage(s *dg.Session, xp *xp.System) Assigner {
	return &roleSystem{
		dgs:   s,
		xp:    xp,
		tiers: &autoRanks{Tiers: make(map[string]float64)},
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
	case "lar": // list auto ranks
		return r.listTiers()
	case "aar": // add auto rank
		ownerID, found := os.LookupEnv("BOT_OWNER")
		if !found {
			return nil, fmt.Errorf(
				"XP Execute failed: \"BOT_OWNER\" environment variable not found",
			)
		}
		if msg.Author.ID == ownerID {
			return r.addAutoRank(x.AutoRankFile, req[2], req[3])
		}
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

func (r *roleSystem) addAutoRank(file, roleName, minXP string) (MsgEmbed, error) {
	// Command: =xp aar roleName minXP
	xp, err := strconv.ParseFloat(minXP, 64)
	if err != nil {
		return nil, fmt.Errorf("addAutoRanks() return: %s", err)
	}
	r.tiers.Tiers[roleName] = xp

	err = r.saveAutoRanks(file)
	if err != nil {
		return nil, fmt.Errorf("saveAutoRanks failed: %s", err)
	}

	return &dg.MessageEmbed{
		Description: fmt.Sprintf("Added %s to be awarded at >= %.2f", roleName, xp)}, nil
}

// LoadAutoRanks loads the saved xp data from the json file
func (r *roleSystem) LoadAutoRanks(file string) error {
	savedTiers, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	err = json.Unmarshal(savedTiers, r.tiers)
	if err != nil {
		return err
	}
	return nil
}

// SaveXP saves the current struct data to a json file
func (r *roleSystem) saveAutoRanks(file string) error {
	json, err := json.MarshalIndent(r.tiers, "", "  ")
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

// AutoPromote checks the user's current XP after each messag sent
// and promotes the user to the correct role if crosses a certain threshold.
func (r *roleSystem) AutoPromote(msg *dg.MessageCreate) error {
	userID := msg.Author.ID
	guildID := msg.GuildID

	// If it's a bot; skip promotion
	if msg.Author.Bot {
		return nil
	}

	// Get roles set in the guild (server)
	roles, err := r.dgs.GuildRoles(guildID)
	if err != nil {
		return fmt.Errorf("GuildRoles failed: %s", err)
	}

	// Set up the map of role names to their IDs
	roleMap := make(map[string]string)
	for _, role := range roles {
		roleMap[role.Name] = role.ID
	}

	// Get user's current total xp
	totalXP, ok := r.xp.Data.Users[userID]
	if !ok {
		return fmt.Errorf("user ID (%s) not found", userID)
	}

	// Set all roles that the user's xp allows
	var roleID string
	for roleName, minXP := range r.tiers.Tiers {
		if totalXP >= minXP {
			roleID = roleMap[roleName]
			err = r.dgs.GuildMemberRoleAdd(guildID, userID, roleID)
			if err != nil {
				return fmt.Errorf("GuildMemberRoleAdd failed: %s", err)
			}
		}
	}

	return nil
}

func (r *roleSystem) listTiers() (MsgEmbed, error) {
	flipTiers := make(map[float64]string)
	rankOrder := []float64{}

	// create map of min xp and tier names
	for name, threshold := range r.tiers.Tiers {
		flipTiers[threshold] = name
		rankOrder = append(rankOrder, threshold)
	}

	// sort  the slice of min xp into ascending order
	sortedRanks := sort.Float64Slice(rankOrder)

	// Create output string in the order of sortedRanks
	var tierList string
	for _, xp := range sortedRanks {
		tierList = fmt.Sprintf("%s\n%s: %.2f", tierList, flipTiers[xp], xp)
	}
	return &dg.MessageEmbed{Title: "Auto Rank XP", Description: tierList}, nil
}
