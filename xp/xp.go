package xp

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"sync"
	"time"

	p "git.sr.ht/~jrswab/akordo/plugins"
	dg "github.com/bwmarrin/discordgo"
)

// DefaultFile is the path where the xp data is saved.
const DefaultFile string = "xp.json"
const messagePoints float64 = 0.01

// Exp is the interface for interacting with the xp methods
type Exp interface {
	LoadXP(file string) error
	ManipulateXP(action string, msg *dg.MessageCreate)
	AutoSaveXP()
	ReturnXp(req []string, msg *dg.MessageCreate) (string, error)
}

// System holds all data needed to execute the functionality.
type System struct {
	data    *xpData
	callRec *p.Record
	mutex   *sync.Mutex
}

// DataStore holds the experience gained by each user.
type xpData struct {
	Users map[string]float64 `json:"users"`
}

// NewXpStore creates the experience data storage map for the session
func NewXpStore(mtx *sync.Mutex) Exp {
	return &System{
		data:    &xpData{Users: make(map[string]float64)},
		callRec: p.NewRecorder(),
		mutex:   mtx,
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

	switch action {
	case "addMessagePoints":
		x.awardActivity(msg)
	case "save":
		x.saveXP(DefaultFile)
	}

	x.mutex.Unlock()
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

// ReturnXp checks the user's request and returns xp data based on the command entered.
func (x *System) ReturnXp(req []string, msg *dg.MessageCreate) (string, error) {
	if len(req) < 2 {
		return x.userXp(msg)
	}

	return "Sorry, the developers have not added that feature yet :sob:", nil
}

func (x *System) userXp(msg *dg.MessageCreate) (string, error) {
	alertUser, tooSoon := x.callRec.CheckLastAsk(msg)
	if tooSoon {
		return alertUser, nil
	}

	xpFloat, ok := x.data.Users[msg.Author.ID]
	if !ok {
		return "", fmt.Errorf("%s, you have not earned any XP", msg.Author.Username)
	}
	xp := strconv.FormatFloat(xpFloat, 'f', 2, 64)
	return xp, nil
}
