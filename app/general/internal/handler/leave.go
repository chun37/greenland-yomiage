package handler

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func (h *Handler) Leave(s *discordgo.Session, i *discordgo.InteractionCreate) {
	vc, ok := s.VoiceConnections[i.GuildID]
	if !ok {
		log.Println("voice connection is not found")
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "参加中のボイスチャンネルはありません。",
			},
		})
		return
	}

	if err := vc.Disconnect(); err != nil {
		log.Println("failed to disconnect voice connection:", err)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "ボイスチャンネルからの退出に失敗しました。",
			},
		})
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "ばいばい〜",
		},
	})
}
