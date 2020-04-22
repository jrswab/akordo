package xp

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"

	dg "github.com/bwmarrin/discordgo"
)

const messagePoints float64 = 0.01

// GamePoints is exported to be used when games need to tell the winner how many
// points they have received.
const GamePoints float64 = 10

// Exp is the interface for interacting with the xp methods
type Exp interface {
	LoadXP()
	ManipulateXP(action string, msg *dg.MessageCreate)
	AutoSaveXP()
}

// System holds all data needed to execute the functionality.
type System struct {
	data   *xpData
	mutex  *sync.Mutex
	User   string
	Points float64
	Award  float64
}

// DataStore holds the experience gained by each user.
type xpData struct {
	Users map[string]float64 `json:"users"`
}

// NewXpStore creates the experience data storage map for the session
func NewXpStore(mtx *sync.Mutex) Exp {
	return &System{
		data:  &xpData{Users: make(map[string]float64)},
		mutex: mtx,
	}
}

// LoadXP loads the saved xp data from the json file
func (x *System) LoadXP() {
	savedXp, err := ioutil.ReadFile("xp.json")
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(savedXp, x.data)
	if err != nil {
		log.Fatal(err)
	}
}

// ManipulateXP is used by any part of the program that needs to read or write
// data after startup.
func (x *System) ManipulateXP(action string, msg *dg.MessageCreate) {
	x.mutex.Lock()

	switch action {
	case "addMessagePoints":
		x.awardActivity(msg)
	case "save":
		x.saveXP()
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

func (x *System) gameReward(msg *dg.MessageCreate) {
	award := 1.00 // No reward bonus
	user := msg.Author.ID

	// Don't award points to the bot
	// Set `BOT_ID` as an environment variable to exclude the bot.
	if user == checkBotID() {
		return
	}

	x.writeToXpMap(user, award, GamePoints)
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
func (x *System) saveXP() {
	json, err := json.MarshalIndent(x.data, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	// Write to data to a file
	err = ioutil.WriteFile("xp.json", json, 0600)
	if err != nil {
		log.Fatal(err)
	}
}

func checkBotID() string {
	// Ignoring second value in case the bot owner wants to allow
	// the bot to gain experience.
	botID, _ := os.LookupEnv("BOT_TOKEN")
	return botID
}
