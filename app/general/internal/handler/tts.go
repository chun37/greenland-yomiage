package handler

import (
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/chun37/greenland-yomiage/general/internal/speaker"
)

func (h *Handler) TTS(messages chan speaker.SpeechMessage, x chan struct{}) func(s *discordgo.Session, m *discordgo.MessageCreate) {
	return func(s *discordgo.Session, m *discordgo.MessageCreate) {
		// m.Author.ID == s.State.User.ID: 自分のメッセージ
		// m.Author.Bot: Bot のメッセージ
		// h.props.Config.TargetChannelID != m.ChannelID: 読み上げチャンネル以外
		if m.Author.ID == s.State.User.ID || m.Author.Bot || h.props.Config.TargetChannelID != m.ChannelID {
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

		// vs == nil: VC に参加してない
		// !vs.SelfMute: ミュートしていない
		// vs.SelfDeaf: スピーカーミュートしている
		// vs.Mute: サーバーミュートされている
		/*if vs == nil || !vs.SelfMute || vs.SelfDeaf || vs.Mute {
			return
		}*/

		v, err := h.joinvc(s, vs.GuildID, vs.ChannelID)
		if err != nil {
			log.Println("failed to join voice channel:", err)
			return
		}

		time.Sleep(time.Millisecond * 500)

		messages <- speaker.SpeechMessage{v, m.Content}
	}
}
