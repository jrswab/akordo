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
	sess, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalf("Session creation error: %s", err)
	}

	sess.AddHandler(controller.ReceiveMessage)

	if err = sess.Open(); err != nil {
		log.Fatalf("Open session error: %s", err)
		return
	}

	// Wait here until until told to terminate. (ie: ctrl+c)
	fmt.Println("Bot is now running.\n\nPress CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	if err := sess.Close(); err != nil {
		log.Fatalf("session.Close failed: %s", err)
	}
}
