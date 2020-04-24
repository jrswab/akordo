package xp

import (
	"fmt"
	"regexp"
	"strings"

	dg "github.com/bwmarrin/discordgo"
)

// Execute is the method used to run the correct method based on user input.
func (x *System) Execute(req []string, msg *dg.MessageCreate) (string, error) {
	// When user runs the xp command with alone return that user's XP
	if len(req) < 2 {
		return x.userXp("", "", msg)
	}

	switch req[1] {
	case "save":
		// DefaultFile is declared in xp/xp.go
		if err := x.saveXP(DefaultFile); err != nil {
			return "", err
		}
		return "XP data saved!", nil
	default:
		return x.returnXp(req, msg)
	}
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
