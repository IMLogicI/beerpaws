package bot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
)

func (b *Bot) getPointsByDiscordIDHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	count, err := b.getPointsByDiscordID(m.Author.ID)
	if err != nil {
		_, _ = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Что-то пошло не так! : %v", err))
		return
	}

	_, _ = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("У вас на баллансе %d баллов", count))
}

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
