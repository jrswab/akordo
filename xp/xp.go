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

// Exp is the interface for interacting with the xp methods
type Exp interface {
	AwardActivity(mutex *sync.Mutex, msg *dg.MessageCreate)
	LoadXP()
	SaveXP(mutex *sync.Mutex)
	AutoSaveXP(mutex *sync.Mutex)
}

// DataStore holds the experience gained by each user.
type xpStore struct {
	Users map[string]float64 `json:"users"`
}

// NewXpStore creates the experience data storage map for the session
func NewXpStore() Exp {
	return &xpStore{
		Users: make(map[string]float64),
	}
}

// LoadXP loads the saved xp data from the json file
func (x *xpStore) LoadXP() {
	savedXp, err := ioutil.ReadFile("xp.json")
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(savedXp, x)
	if err != nil {
		log.Fatal(err)
	}
}

// SaveXP saves the current struct data to a json file
func (x *xpStore) SaveXP(mutex *sync.Mutex) {
	// Lock the map to avoid race conditions
	mutex.Lock()
	defer mutex.Unlock()

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

func (x *xpStore) AutoSaveXP(mutex *sync.Mutex) {
	for true {
		select {
		case <-time.After(5 * time.Minute):
			x.SaveXP(mutex)
		}
	}
}

// AwardXP stores the earned experience into the DataStore struct.
func (x *xpStore) AwardActivity(mutex *sync.Mutex, msg *dg.MessageCreate) {
	award := len(msg.Content)
	user := msg.Author.ID
	points := 0.01 // per character points

	// Don't award points to the bot
	// Set `BOT_ID` as an environment variable to exclude the bot.
	if user == checkBotID() {
		return
	}

	x.writeToXpMap(mutex, user, float64(award), points)
}

func (x *xpStore) writeToXpMap(mutex *sync.Mutex, user string, award, points float64) {
	mutex.Lock()
	defer mutex.Unlock()

	if xp, ok := x.Users[user]; ok {
		x.Users[user] = xp + (award * points)
		return
	}

	x.Users[user] = float64(award) * points
}

func checkBotID() string {
	// Ignoring second value in case the bot owner wants to allow
	// the bot to gain experience.
	botID, _ := os.LookupEnv("BOT_TOKEN")
	return botID
}
