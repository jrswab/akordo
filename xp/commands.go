package xp

import (
	"fmt"
	"log"
	"regexp"
	"sort"
	"strings"

	dg "github.com/bwmarrin/discordgo"
)

// AutoRankFile is the default file for loading and saving auto promote ranks
const AutoRankFile string = "data/autoRanks.json"

// MsgEmbed is used to shorten the name of the original embed type from discordGo
type MsgEmbed *dg.MessageEmbed

// Execute is the method used to run the correct method based on user input.
func (x *System) Execute(req []string, msg *dg.MessageCreate) (MsgEmbed, error) {
	// When user runs the xp command with alone return that user's XP
	if len(req) < 2 {
		return x.returnXp(req, msg)
	}

	switch req[1] {
	case "save":
		// DefaultFile is declared in xp/xp.go
		if err := x.saveXP(XpFile); err != nil {
			return nil, err
		}
		return &dg.MessageEmbed{Description: "XP data saved!"}, nil
	case "lb":
		return x.leaderBoard(msg)
	}
	return x.returnXp(req, msg)
}

// ReturnXp checks the user's request and returns xp data based on the command entered.
func (x *System) returnXp(req []string, msg *dg.MessageCreate) (MsgEmbed, error) {
	// When user runs the xp command with alone return that user's XP
	if len(req) < 2 {
		return x.userXp(msg.Author.Username, msg.Author.ID, msg)
	}

	var id string
	regEx := fmt.Sprintf("(?m)^<@!\\w+>")
	var re = regexp.MustCompile(regEx)
	match := re.MatchString(req[1])
	if match {
		id = strings.TrimPrefix(req[1], "<@!")
		id = strings.TrimSuffix(id, ">")

		member, err := x.dgs.GuildMember(msg.GuildID, id)
		if err != nil {
			return nil, fmt.Errorf("returnXP() call to GuildMember returned: %s", err)
		}

		return x.userXp(member.User.Username, id, msg)
	}

	return &dg.MessageEmbed{Description: "User not found... Did you use `@`?"}, nil
}

func (x *System) userXp(name, userID string, msg *dg.MessageCreate) (MsgEmbed, error) {
	alertUser, tooSoon := x.callRec.CheckLastAsk(msg)
	if tooSoon {
		return &dg.MessageEmbed{Description: alertUser}, nil
	}

	xp, ok := x.Data.Users[userID]
	if !ok {
		return &dg.MessageEmbed{Description: fmt.Sprintf("%s has not earned any XP", name)}, nil
	}

	return &dg.MessageEmbed{Description: fmt.Sprintf("%s has a total of %.2f xp", name, xp)}, nil
}

// Currently not used; needs a way to look up nicknames before checking the actual
// usernames to avoid returning incorrect data.
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

func (x *System) leaderBoard(msg *dg.MessageCreate) (MsgEmbed, error) {
	flippedMap := make(map[float64]string)
	flipSlice := []float64{}
	for key, value := range x.Data.Users {
		flipSlice = append(flipSlice, value)
		flippedMap[value] = key
	}

	sort.Float64s(flipSlice)

	totalUsers := len(flipSlice) - 1
	topStop := totalUsers - 10
	var top10 string

	rank := 1
	for i := totalUsers; i > topStop; i-- {
		userID := flippedMap[flipSlice[i]]
		user, err := x.dgs.GuildMember(msg.GuildID, userID)
		if err != nil {
			log.Printf("leaderboard() GuildMember() returned an error: %s", err)
			continue
		}
		top10 = fmt.Sprintf("%s\n%d) %s (%.2f)", top10, rank, user.User.Username, flipSlice[i])
		rank++
	}

	embed := &dg.MessageEmbed{
		Title:       "Top 10",
		Description: top10,
	}

	return embed, nil
}
