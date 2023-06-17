package handler

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func (h *Handler) Disconnect(s *discordgo.Session, u *discordgo.VoiceStateUpdate) {
	vc, ok := s.VoiceConnections[u.GuildID]
	if !ok {
		return
	}
	g, err := s.State.Guild(vc.GuildID)
	if err != nil {
		log.Println("failed to get guild from guild id:", err)
		return
	}
	members := func() []*discordgo.Member {
		ms := make([]*discordgo.Member, 0)
		for _, vs := range g.VoiceStates {
			if vs.ChannelID != vc.ChannelID {
				continue
			}
			m, err := s.State.Member(vs.GuildID, vs.UserID)
			if err != nil {
				continue
			}
			ms = append(ms, m)
		}
		return ms
	}()

	if len(members) >= 2 {
		return
	}

	if err := vc.Disconnect(); err != nil {
		log.Println("failed to disconnect voice connection:", err)
	}
}
