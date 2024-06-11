package bot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
)

func (b *Bot) getPointsByDiscordID(
	discordID string,
) (int64, error) {
	user, err := b.userService.GetUserByDiscordID(discordID)
	if err != nil {
		return 0, err
	}

	if user == nil {
		return 0, fmt.Errorf("пользователь не зарегистрирован в системе")
	}

	return b.pointService.GetPointsByUserID(user.ID)
}

func getRulesInteraction(s *discordgo.Session, i *discordgo.InteractionCreate, b *Bot) {
	count, err := b.getPointsByDiscordID(i.Interaction.Member.User.ID)
	if err != nil {
		errInteraction(s, i)
	}
	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:   discordgo.MessageFlagsEphemeral,
			Content: fmt.Sprintf("У вас на счету %v баллов", count),
			Title:   "Баллы!",
		},
	})
}
