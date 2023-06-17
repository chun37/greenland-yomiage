package handler

import (
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/chun37/greenland-yomiage/general/internal/speaker"
)

func (h *Handler) TTS(messages chan speaker.SpeechMessage) func(s *discordgo.Session, m *discordgo.MessageCreate) {
	return func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return
		}

		if m.Author.Bot {
			return
		}

		if h.props.Config.TargetChannelID != m.ChannelID {
			return
		}

		guild, err := s.State.Guild(m.GuildID)
		if err != nil {
			log.Println("failed to get guild:", err)
			return
		}

		vs := func() *discordgo.VoiceState {
			for _, state := range guild.VoiceStates {
				if state.UserID == m.Author.ID {
					return state
				}
			}
			return nil
		}()
		if vs == nil {
			return
		}

		v, err := s.ChannelVoiceJoin(vs.GuildID, vs.ChannelID, false, true)
		if err != nil {
			log.Println("failed to join voice channel:", err)
			return
		}

		messages <- speaker.SpeechMessage{v, m.Content}
	}
}
