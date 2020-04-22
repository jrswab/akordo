package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"git.sr.ht/~jrswab/akordo/controller"
	"github.com/bwmarrin/discordgo"
)

var token string

func init() {
	flag.StringVar(&token, "t", "", "Bot Token")
	flag.Parse()
}

func main() {
	// If token is not passed by -t check for and environment variable
	if token == "" {
		var found bool
		token, found = os.LookupEnv("BOT_TOKEN")
		if !found {
			log.Fatalf("Please pass in your bot token with -t or set the \"BOT_TOKEN\" environment variable.")
		}
	}

	sess, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalf("Session creation error: %s", err)
	}

	// Create a the custom controller to pass data to ReceiveMessage and the plugins
	sd := controller.NewSessionData(sess)
	// Load saved XP data into the struct created by NewSessionData
	if _, err := os.Stat("xp.json"); err == nil {
		sd.UserXP.LoadXP()
	}

	// start the Goroutine to automatically save earned XP
	go sd.UserXP.AutoSaveXP()

	// Watch for new messages
	sess.AddHandler(sd.NewMessage)

	if err = sess.Open(); err != nil {
		log.Fatalf("Open session error: %s", err)
		return
	}

	// Wait here until until told to terminate. (ie: ctrl+c)
	fmt.Println("Bot is now running.\n\nPress CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Save data
	sd.UserXP.ManipulateXP("save", &discordgo.MessageCreate{})

	// Close the session
	if err := sess.Close(); err != nil {
		log.Fatalf("session.Close failed: %s", err)
	}
}
