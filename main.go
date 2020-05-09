package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"gitlab.com/technonauts/akordo/controller"
	"gitlab.com/technonauts/akordo/load"
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

	load.SavedData(sd)

	// start the Goroutine to automatically save earned XP
	go sd.XP.AutoSaveXP()

	// Watch for new messages
	sess.AddHandler(sd.NewMessage)

	if err = sess.Open(); err != nil {
		log.Fatalf("Open session error: %s", err)
		return
	}

	// Wait here until until told to terminate. (ie: ctrl+c)
	log.Println("Bot is now running.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Save data on exit
	sd.XP.ManipulateXP("save", &discordgo.MessageCreate{})

	// Close the session
	if err := sess.Close(); err != nil {
		log.Fatalf("session.Close failed: %s", err)
	}
}
