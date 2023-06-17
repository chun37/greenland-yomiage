package handler

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/chun37/greenland-yomiage/internal/usecase/tts"
)

func (h *Handler) TTS(s *discordgo.Session, m *discordgo.MessageCreate) {
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
		log.Println("failed to get voice state:", err)
		return
	}

	v, err := s.ChannelVoiceJoin(vs.GuildID, vs.ChannelID, false, true)
	if err != nil {
		log.Println("failed to join voice channel:", err)
		return
	}
	err = v.Speaking(true)
	if err != nil {
		fmt.Println("Couldn't set speaking", err)
		return
	}

	done := make(chan struct{})
	opusChunks := make(chan []byte, 3)
	defer close(done)
	defer close(opusChunks)
	h.props.TTSUsecase.Do(tts.UsecaseParam{
		Text:       m.Content,
		OpusChunks: opusChunks,
		Done:       done,
	})

	// Send not "speaking" packet over the websocket when we finish
	defer func() {
		err := v.Speaking(false)
		if err != nil {
			fmt.Println("Couldn't stop speaking", err)
		}
	}()

	for {
		select {
		case opus := <-opusChunks:
			v.OpusSend <- opus
		case <-done:
			return
		}
	}
}
