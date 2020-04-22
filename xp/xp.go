package xp

import (
	"encoding/json"
	"io/ioutil"
	"log"

	dg "github.com/bwmarrin/discordgo"
)

// Exp is the interface for interacting with the xp methods
type Exp interface {
	AwardXP(msg *dg.MessageCreate)
	LoadXP()
	SaveXP()
}

// DataStore holds the experience gained by each user.
type xpStore struct {
	PointsPerChar float64            `json:"pointsPerChar"`
	Users         map[string]float64 `json:"users"`
}

// NewXpStore creates the experience data storage map for the session
func NewXpStore() Exp {
	return &xpStore{
		PointsPerChar: 0.01,
		Users:         make(map[string]float64),
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
func (x *xpStore) SaveXP() {
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

// AwardXP stores the earned experience into the DataStore struct.
func (x *xpStore) AwardXP(msg *dg.MessageCreate) {
	award := len(msg.Content)
	user := msg.Author.ID

	if xp, ok := x.Users[user]; ok {
		x.Users[user] = xp + (float64(award) * x.PointsPerChar)
		log.Printf("%s - %.2f", user, x.Users[user])
		return
	}

	x.Users[user] = float64(award) * x.PointsPerChar
}
