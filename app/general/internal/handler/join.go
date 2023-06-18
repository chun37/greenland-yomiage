package handler

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

func (h *Handler) Join(s *discordgo.Session, i *discordgo.InteractionCreate) {
	guild, err := s.State.Guild(i.GuildID)
	if err != nil {
		log.Println("failed to get guild:", err)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("ギルド情報の取得に失敗しました。"),
			},
		})
		return
	}

	vs := func() *discordgo.VoiceState {
		for _, state := range guild.VoiceStates {
			if state.UserID == i.Member.User.ID {
				return state
			}
		}
		return nil
	}()

	if vs == nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("ボイスチャンネルに参加してからコマンドを実行してね。"),
			},
		})
		return
	}

	v, err := s.ChannelVoiceJoin(vs.GuildID, vs.ChannelID, false, true)
	if err != nil {
		log.Println("failed to join voice channel:", err)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("ボイスチャンネルへの参加に失敗しました。"),
			},
		})
		return
	}

	c, err := s.State.Channel(v.ChannelID)
	if err != nil {
		log.Println("failed to get channel:", err)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("チャンネル情報の取得に失敗しました。"),
			},
		})
		return
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("[%s]に参加したよ", c.Name),
		},
	})
}
