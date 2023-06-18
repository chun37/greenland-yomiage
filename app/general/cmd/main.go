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
	"github.com/chun37/greenland-yomiage/general/internal/speaker"
)

// Variables used for command line parameters
var (
	Token   string
	GuildID string
)

func init() {
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.StringVar(&GuildID, "g", "", "Slash Commands Guild")
	flag.Parse()
}

func main() {
	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		log.Fatalf("error creating Discord session,", err)
	}

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.IntentsAll

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		log.Fatalf("error opening connection,", err)
	}

	// Register the messageCreate func as a callback for MessageCreate events.

	cfg := config.Config{
		TargetChannelID: "773094074269958154",
	}
	externalDeps := initialize.NewExternalDependencies()
	usecases := initialize.NewUsecases(externalDeps)
	hp := initialize.NewHandlerProps(cfg)

	messages := make(chan speaker.SpeechMessage, 10)

	hdr := handler.New(hp, messages)
	dg.AddHandler(hdr.TTS(messages))
	dg.AddHandler(hdr.Disconnect)

	interactionHandler, slashCommandIDs := hdr.Interaction(dg, GuildID)
	dg.AddHandler(interactionHandler)

	spkr := speaker.NewSpeaker(usecases.TTSUsecase, messages)
	go spkr.Run()

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	log.Println("Removing commands...")
	// // We need to fetch the commands, since deleting requires the command ID.
	// // We are doing this from the returned commands on line 375, because using
	// // this will delete all the commands, which might not be desirable, so we
	// // are deleting only the commands that we added.
	// registeredCommands, err := s.ApplicationCommands(s.State.User.ID, *GuildID)
	// if err != nil {
	// 	log.Fatalf("Could not fetch registered commands: %v", err)
	// }

	for _, scid := range slashCommandIDs {
		err := dg.ApplicationCommandDelete(dg.State.User.ID, GuildID, scid)
		if err != nil {
			log.Panicf("Cannot delete command: %+v", err)
		}
	}

	log.Println("Gracefully shutting down.")

	// Cleanly close down the Discord session.
	dg.Close()
}
