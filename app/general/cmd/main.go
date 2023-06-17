package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/chun37/greenland-yomiage/general/internal/handler"
)

// Variables used for command line parameters
var (
	Token string
)
var buffer = make([][]byte, 0)

func init() {
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
	err := loadSound()
	if err != nil {
		return
	}
}

func main() {

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		log.Fatalf("error creating Discord session,", err)
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(handler.TTS)
	dg.AddHandler(handler.Play)
	dg.AddHandler(handler.AirHorn(buffer))

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.IntentsAll

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		log.Fatalf("error opening connection,", err)
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}
