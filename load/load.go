package load

import (
	"log"
	"os"

	"gitlab.com/technonauts/akordo/controller"
	"gitlab.com/technonauts/akordo/plugins"
	"gitlab.com/technonauts/akordo/roles"
	"gitlab.com/technonauts/akordo/xp"
)

// SavedData loads all saved dato to the correct structs for the specified commands.
func SavedData(sd *controller.SessionData) {
	// Load saved XP data into the struct created by NewSessionData
	if _, err := os.Stat(xp.XpFile); err == nil {
		if err := sd.XP.LoadXP(xp.XpFile); err != nil {
			log.Fatalf("error loading xp data file: %s", err)
		}
	}
	// Load saved role data into the struct created by NewSessionData
	if _, err := os.Stat(xp.AutoRankFile); err == nil {
		if err := sd.Roles.LoadAutoRanks(xp.AutoRankFile); err != nil {
			log.Fatalf("error loading role file: %s", err)
		}
	}
	// Load saved self assign role data into the struct created by NewSessionData
	if _, err := os.Stat(roles.SelfAssignFile); err == nil {
		if err := sd.Roles.LoadSelfAssignRoles(roles.SelfAssignFile); err != nil {
			log.Fatalf("error loading self assign role file: %s", err)
		}
	}
	// Load saved banned word data into the struct created by NewSessionData
	if _, err := os.Stat(plugins.BannedWordsPath); err == nil {
		if err := sd.Blacklist.LoadBannedWordList(plugins.BannedWordsPath); err != nil {
			log.Fatalf("error loading banned words file: %s", err)
		}
	}
	// Load saved base chat role data into the struct created by NewSessionData
	if _, err := os.Stat(plugins.ChatPermissionRole); err == nil {
		if err := sd.Rules.LoadAgreementRole(plugins.ChatPermissionRole); err != nil {
			log.Fatalf("error loading Chat Permission file: %s", err)
		}
	}
}
