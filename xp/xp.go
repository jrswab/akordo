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
const gamePoints float64 = 10

// Exp is the interface for interacting with the xp methods
type Exp interface {
	LoadXP()
	ManipulateXP(action string)
	AutoSaveXP(mutex *sync.Mutex)
	SetXpMsg(msg *dg.MessageCreate)
}

type xpSystem struct {
	data   *xpData
	msg    *dg.MessageCreate
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
	return &xpSystem{
		data:  &xpData{Users: make(map[string]float64)},
		mutex: mtx,
	}
}

func (x *xpSystem) SetXpMsg(msg *dg.MessageCreate) {
	x.msg = msg
}

// LoadXP loads the saved xp data from the json file
func (x *xpSystem) LoadXP() {
	savedXp, err := ioutil.ReadFile("xp.json")
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(savedXp, x)
	if err != nil {
		log.Fatal(err)
	}
}

func (x *xpSystem) ManipulateXP(action string) {
	x.mutex.Lock()
	defer x.mutex.Unlock()

	switch action {
	case "addMessagePoints":
		x.awardActivity(x.msg)
	case "addGamePoints":
		x.gameReward(x.msg)
	case "save":
		x.saveXP()
	}

}

// AwardXP stores the earned experience into the DataStore struct.
func (x *xpSystem) awardActivity(msg *dg.MessageCreate) {
	award := len(msg.Content)
	user := msg.Author.ID

	// Don't award points to the bot
	// Set `BOT_ID` as an environment variable to exclude the bot.
	if user == checkBotID() {
		return
	}

	x.writeToXpMap(user, float64(award), messagePoints)
}

func (x *xpSystem) gameReward(msg *dg.MessageCreate) {
	award := 1.00 // No reward bonus
	user := msg.Author.ID

	// Don't award points to the bot
	// Set `BOT_ID` as an environment variable to exclude the bot.
	if user == checkBotID() {
		return
	}

	x.writeToXpMap(user, award, gamePoints)
}

func (x *xpSystem) writeToXpMap(user string, award, points float64) {
	if xp, ok := x.data.Users[user]; ok {
		x.data.Users[user] = xp + (award * points)
		return
	}

	x.data.Users[user] = float64(award) * points
}

func (x *xpSystem) AutoSaveXP(mutex *sync.Mutex) {
	for true {
		select {
		case <-time.After(5 * time.Minute):
			x.ManipulateXP("save")
		}
	}
}

// SaveXP saves the current struct data to a json file
func (x *xpSystem) saveXP() {
	// Lock the map to avoid race conditions
	x.mutex.Lock()
	defer x.mutex.Unlock()

	json, err := json.MarshalIndent(x, "", "  ")
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
