package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	config "github.com/chun37/greenland-yomiage/general/internal/config"
	"github.com/chun37/greenland-yomiage/general/internal/handler"
	"github.com/chun37/greenland-yomiage/general/internal/initialize"
)

// Variables used for command line parameters
var (
	Token string
)

func init() {
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

func main() {
	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		log.Fatalf("error creating Discord session,", err)
	}

	// Register the messageCreate func as a callback for MessageCreate events.

	cfg := config.Config{
		TargetChannelID: "773181736988573697",
	}
	usecases := initialize.NewUsecases()
	hp := initialize.NewHandlerProps(usecases, cfg)
	hdr := handler.NewHandler(hp)
	dg.AddHandler(hdr.TTS)
	dg.AddHandler(hdr.Play)

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
