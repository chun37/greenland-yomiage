package speaker

import "github.com/bwmarrin/discordgo"

type SpeechMessage struct {
	VoiceConnection *discordgo.VoiceConnection
	Text            string
}
