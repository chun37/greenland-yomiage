package handler

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"golang.org/x/xerrors"
)

func (h *Handler) AddWord(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options

	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	word, err := func() (string, error) {
		opt, ok := optionMap["word"]
		if !ok {
			return "", xerrors.New("failed to get word option")
		}
		return opt.StringValue(), nil
	}()
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "単語が指定されていません。",
			},
		})
		return
	}

	pronunciation, err := func() (string, error) {
		opt, ok := optionMap["pronunciation"]
		if !ok {
			return "", xerrors.New("failed to get pronunciation option")
		}
		return opt.StringValue(), nil
	}()
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "読みが指定されていません。",
			},
		})
		return
	}

	accentType, err := func() (int, error) {
		opt, ok := optionMap["accent_type"]
		if !ok {
			return 0, xerrors.New("failed to get accent_type option")
		}
		return int(opt.IntValue()), nil
	}()
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "アクセント核位置が指定されていません。",
			},
		})
		return
	}

	if err := h.props.DictionaryAddUsecase.Do(word, pronunciation, accentType); err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			// Ignore type for now, they will be discussed in "responses"
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("辞書への単語登録に失敗しました。\n%s", err),
			},
		})

		return
	}
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		// Ignore type for now, they will be discussed in "responses"
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("辞書への単語登録に成功しました。\n[%s](%s)", word, pronunciation),
		},
	})
}
