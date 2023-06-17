package handler

import (
	"github.com/bwmarrin/discordgo"
)

func (h *Handler) TTS(s *discordgo.Session, m *discordgo.MessageCreate) {
	h.props.TTSUsecase.Do(m.Content)
}
