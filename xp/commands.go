package xp

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"

	dg "github.com/bwmarrin/discordgo"
)

// AutoRankFile is the default file for loading and saving auto promote ranks
const AutoRankFile string = "data/autoRanks.json"

// Execute is the method used to run the correct method based on user input.
func (x *System) Execute(req []string, msg *dg.MessageCreate) (string, error) {
	// When user runs the xp command with alone return that user's XP
	if len(req) < 2 {
		return x.userXp("", "", msg)
	}

	switch req[1] {
	case "save":
		// DefaultFile is declared in xp/xp.go
		if err := x.saveXP(XpFile); err != nil {
			return "", err
		}
		return "XP data saved!", nil
	case "aar": // add auto rank
		ownerID, found := os.LookupEnv("BOT_OWNER")
		if !found {
			return "", fmt.Errorf(
				"XP Execute failed: \"BOT_OWNER\" environment variable not found",
			)
		}
		if msg.Author.ID == ownerID {
			return x.addAutoRank(AutoRankFile, req[2], req[3])
		}
	}
	return x.returnXp(req, msg)
}

// ReturnXp checks the user's request and returns xp data based on the command entered.
func (x *System) returnXp(req []string, msg *dg.MessageCreate) (string, error) {
	var id, user string
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

type autoRanks struct {
	Tiers map[string]float64 `json:"tiers"` // map of total xp of role IDs
}

func (x *System) addAutoRank(file, roleName, minXP string) (string, error) {
	// Command: =xp aar roleName minXP
	xp, err := strconv.ParseFloat(minXP, 64)
	if err != nil {
		return "", fmt.Errorf("addAutoRanks() return: %s", err)
	}
	x.tiers.Tiers[roleName] = xp

	err = x.saveAutoRanks(file)
	if err != nil {
		return "", fmt.Errorf("saveAutoRanks failed: %s", err)
	}

	return fmt.Sprintf("Added %s to be awarded at >= %.2f", roleName, xp), nil
}

// LoadAutoRanks loads the saved xp data from the json file
func (x *System) LoadAutoRanks(file string) error {
	savedTiers, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	err = json.Unmarshal(savedTiers, x.tiers)
	if err != nil {
		return err
	}
	return nil
}

// SaveXP saves the current struct data to a json file
func (x *System) saveAutoRanks(file string) error {
	json, err := json.MarshalIndent(x.tiers, "", "  ")
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
