package handler

import "github.com/bwmarrin/discordgo"

func (h *Handler) Interaction(dg *discordgo.Session, guildID string) (func(s *discordgo.Session, i *discordgo.InteractionCreate), []string) {
	basicCommand := &discordgo.ApplicationCommand{
		Name:        "basic-command",
		Description: "Basic command",
	}

	createdBasicCommand, err := dg.ApplicationCommandCreate(dg.State.User.ID, guildID, basicCommand)
	if err != nil {

	}

	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.ApplicationCommandData().Name {
		case "basic-command":
			h.BasicCommand(s, i)
		}
	}, []string{createdBasicCommand.ID}
}
