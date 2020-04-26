package xp

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"time"

	p "git.sr.ht/~jrswab/akordo/plugins"
	dg "github.com/bwmarrin/discordgo"
)

// XpFile is the path where the xp data is saved.
const XpFile string = "xp.json"
const messagePoints float64 = 0.01

// AkSession allows for tests to mock discordgo session method calls
type AkSession interface {
	GuildMember(guildID string, userID string) (st *dg.Member, err error)
	GuildMembers(guildID string, after string, limit int) (st []*dg.Member, err error)
	GuildRoles(guildID string) (st []*dg.Role, err error)
	GuildMemberRoleAdd(guildID string, userID string, roleID string) (err error)
}

// Exp is the interface for interacting with the xp methods
type Exp interface {
	LoadXP(file string) error
	ManipulateXP(action string, msg *dg.MessageCreate)
	AutoSaveXP()
	Execute(req []string, msg *dg.MessageCreate) (string, error)
	AutoPromote(msg *dg.MessageCreate) error
	LoadAutoRanks(file string) error
}

// System holds all data needed to execute the functionality.
type System struct {
	data    *xpData
	tiers   *autoRanks
	callRec *p.Record
	mutex   *sync.Mutex
	dgs     AkSession
}

// DataStore holds the experience gained by each user.
type xpData struct {
	Users map[string]float64 `json:"users"`
}

// NewXpStore creates the experience data storage map for the session
func NewXpStore(mtx *sync.Mutex, s *dg.Session) Exp {
	return &System{
		data:    &xpData{Users: make(map[string]float64)},
		tiers:   &autoRanks{Tiers: make(map[string]float64)},
		callRec: p.NewRecorder(),
		mutex:   mtx,
		dgs:     s,
	}
}

// LoadXP loads the saved xp data from the json file
func (x *System) LoadXP(file string) error {
	savedXp, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	err = json.Unmarshal(savedXp, x.data)
	if err != nil {
		return err
	}
	return nil
}

// ManipulateXP is used by any part of the program that needs to read or write
// data after startup.
func (x *System) ManipulateXP(action string, msg *dg.MessageCreate) {
	x.mutex.Lock()
	defer x.mutex.Unlock()

	switch action {
	case "addMessagePoints":
		x.awardActivity(msg)
	case "save":
		x.saveXP(XpFile)
	}
}

// AwardXP stores the earned experience into the DataStore struct.
func (x *System) awardActivity(msg *dg.MessageCreate) {
	award := len(msg.Content)
	user := msg.Author.ID

	// Don't award points to the bot
	// Set `BOT_ID` as an environment variable to exclude the bot.
	if user == checkBotID() {
		return
	}

	x.writeToXpMap(user, float64(award), messagePoints)
}

func (x *System) writeToXpMap(user string, award, points float64) {
	if xp, ok := x.data.Users[user]; ok {
		x.data.Users[user] = xp + (award * points)
		return
	}

	x.data.Users[user] = float64(award) * points
}

// AutoSaveXP is launched by main.go before accepting new messages.
// Default time is coded to 5 minutes.
func (x *System) AutoSaveXP() {
	for true {
		select {
		case <-time.After(5 * time.Minute):
			x.ManipulateXP("save", &dg.MessageCreate{})
		}
	}
}

// SaveXP saves the current struct data to a json file
func (x *System) saveXP(file string) error {
	json, err := json.MarshalIndent(x.data, "", "  ")
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

func checkBotID() string {
	// Ignoring second value in case the bot owner wants to allow
	// the bot to gain experience.
	botID, _ := os.LookupEnv("BOT_TOKEN")
	return botID
}

// AutoPromote checks the user's current XP after each messag sent
// and promotes the user to the correct role if crosses a certain threshold.
func (x *System) AutoPromote(msg *dg.MessageCreate) error {
	userID := msg.Author.ID
	guildID := msg.GuildID

	// Get roles set in the guild (server)
	roles, err := x.dgs.GuildRoles(guildID)
	if err != nil {
		return fmt.Errorf("GuildRoles failed: %s", err)
	}

	// Set up the map of role names to their IDs
	roleMap := make(map[string]string)
	for _, role := range roles {
		roleMap[role.Name] = role.ID
	}

	// Get user's current total xp
	totalXP, ok := x.data.Users[userID]
	if !ok {
		return fmt.Errorf("user ID (%s) not found", userID)
	}

	// Set all roles that the user's xp allows
	var roleID string
	for roleName, minXP := range x.tiers.Tiers {
		if totalXP >= minXP {
			roleID = roleMap[roleName]
			break
		}
	}

	if roleID == "" {
		return nil
	}

	err = x.dgs.GuildMemberRoleAdd(guildID, userID, roleID)
	if err != nil {
		return fmt.Errorf("GuildMemberRoleAdd failed: %s", err)
	}

	return nil
}
