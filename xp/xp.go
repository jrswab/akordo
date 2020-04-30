package xp

import (
	"encoding/json"
	"io/ioutil"
	"sync"
	"time"

	p "git.sr.ht/~jrswab/akordo/plugins"
	dg "github.com/bwmarrin/discordgo"
)

// XpFile is the path where the xp data is saved.
const XpFile string = "data/xp.json"
const messagePoints float64 = 0.01

// AkSession allows for tests to mock discordgo session method calls
type AkSession interface {
	GuildMember(guildID string, userID string) (st *dg.Member, err error)
	GuildMembers(guildID string, after string, limit int) (st []*dg.Member, err error)
	GuildRoles(guildID string) (st []*dg.Role, err error)
	GuildMemberRoleAdd(guildID string, userID string, roleID string) (err error)
}

// System holds all data needed to execute the functionality.
type System struct {
	Data    *Data
	callRec *p.Record
	mutex   *sync.Mutex
	dgs     AkSession
}

// Data holds the experience gained by each user.
type Data struct {
	Users map[string]float64 `json:"users"`
}

// NewXpStore creates the experience data storage map for the session
func NewXpStore(mtx *sync.Mutex, s *dg.Session) *System {
	return &System{
		Data:    &Data{Users: make(map[string]float64)},
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

	err = json.Unmarshal(savedXp, x.Data)
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

	// Don't award points to bots
	if msg.Author.Bot {
		return
	}

	x.writeToXpMap(user, float64(award), messagePoints)
}

func (x *System) writeToXpMap(user string, award, points float64) {
	if xp, ok := x.Data.Users[user]; ok {
		x.Data.Users[user] = xp + (award * points)
		return
	}

	x.Data.Users[user] = float64(award) * points
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
	json, err := json.MarshalIndent(x.Data, "", "  ")
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
