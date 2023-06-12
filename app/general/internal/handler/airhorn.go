package handler

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

func AirHorn(buffer [][]byte) func(session *discordgo.Session, create *discordgo.MessageCreate) {
	return func(s *discordgo.Session, m *discordgo.MessageCreate) {

		// Ignore all messages created by the bot itself
		// This isn't required in this specific example but it's a good practice.
		if m.Author.ID == s.State.User.ID {
			return
		}

		// check if the message is "!airhorn"
		if strings.HasPrefix(m.Content, "!airhorn") {

			// Find the channel that the message came from.
			c, err := s.State.Channel(m.ChannelID)
			if err != nil {
				log.Printf("%v", err)
				return
			}

			// Find the guild for that channel.
			g, err := s.State.Guild(c.GuildID)
			if err != nil {
				log.Printf("%v", err)
				return
			}

			// Look for the message sender in that guild's current voice states.
			for _, vs := range g.VoiceStates {
				if vs.UserID == m.Author.ID {
					err = playSound(s, g.ID, vs.ChannelID, buffer)
					if err != nil {
						fmt.Println("Error playing sound:", err)
					}

					return
				}
			}
		}
	}
}

// playSound plays the current buffer to the provided channel.
func playSound(s *discordgo.Session, guildID, channelID string, buffer [][]byte) (err error) {

	// Join the provided voice channel.
	vc, err := s.ChannelVoiceJoin(guildID, channelID, false, true)
	if err != nil {
		return err
	}

	// Sleep for a specified amount of time before playing the sound
	time.Sleep(250 * time.Millisecond)

	// Start speaking.
	vc.Speaking(true)

	// Send the buffer data.
	for _, buff := range buffer {
		vc.OpusSend <- buff
	}

	// Stop speaking
	vc.Speaking(false)

	// Sleep for a specificed amount of time before ending.
	time.Sleep(250 * time.Millisecond)

	// Disconnect from the provided voice channel.
	vc.Disconnect()

	return nil
}
