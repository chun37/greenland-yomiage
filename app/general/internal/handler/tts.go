package handler

import (
	"github.com/bwmarrin/discordgo"
	"github.com/chun37/greenland-yomiage/internal/usecase/tts"
)

func TTS(s *discordgo.Session, m *discordgo.MessageCreate) {
	u := tts.NewUsecase()
	u.Do(m.Content)
}
