package xp

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
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
	dgs     *dg.Session
}

// DataStore holds the experience gained by each user.
type xpData struct {
	Users map[string]float64 `json:"users"`
}

// NewXpStore creates the experience data storage map for the session
func NewXpStore(mtx *sync.Mutex, s *dg.Session) Exp {
	return &System{
		data:    &xpData{Users: make(map[string]float64)},
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
	var id, user string

	// When user runs the xp command with alone return that user's XP
	if len(req) < 2 {
		return x.userXp("", "", msg)
	}

	regEx := fmt.Sprintf("(?m)^<@!\\w+>")
	var re = regexp.MustCompile(regEx)
	match := re.MatchString(req[1])
	if match {
		id = strings.TrimPrefix(req[1], "<@!")
		id = strings.TrimSuffix(id, ">")

		member, err := x.dgs.GuildMember(msg.GuildID, id)
		if err != nil {
			return "", fmt.Errorf("returnXP() call to GuildMember returned: %s", err)
		}

		return x.userXp(member.User.Username, id, msg)
	}

	if !match {
		for idx, word := range req[1:] {
			if idx == 0 {
				user = word
				continue
			}
			user = fmt.Sprintf("%s %s", user, word)
		}
		id, err := x.findUserID(user, msg)
		if err != nil {
			return "", fmt.Errorf("findUserID error: %s", err)
		}
		return x.userXp(user, id, msg)
	}

	return "User not found :thinking:", nil
}

func (x *System) userXp(name, userID string, msg *dg.MessageCreate) (string, error) {
	alertUser, tooSoon := x.callRec.CheckLastAsk(msg)
	if tooSoon {
		return alertUser, nil
	}

	if userID == "" {
		userID = msg.Author.ID
	}

	if name == "" {
		name = msg.Author.Username
	}

	xp, ok := x.data.Users[userID]
	if !ok {
		return fmt.Sprintf("%s has not earned any XP", name), nil
	}

	return fmt.Sprintf("%s has a total of %.2f xp", name, xp), nil
}

func (x *System) findUserID(userName string, msg *dg.MessageCreate) (string, error) {
	var (
		members []*dg.Member
		err     error
	)

	// Get first round of members
	current, err := x.dgs.GuildMembers(msg.GuildID, "", 1000)
	if err != nil {
		return "", err
	}

	for _, names := range current {
		members = append(members, names)
	}

	// if first round has 1000 entries run again until all members are present.
	for len(current) == 1000 {
		lastMember := current[len(current)-1]
		current, err = x.dgs.GuildMembers(msg.GuildID, lastMember.User.ID, 1000)
		if err != nil {
			return "", err
		}

		for _, names := range current {
			members = append(members, names)
		}

	}

	// Create map of usernames and IDs
	userMap := make(map[string]string)
	for _, m := range members {
		userMap[m.User.Username] = m.User.ID
	}

	userID := userMap[userName]

	return userID, nil
}
