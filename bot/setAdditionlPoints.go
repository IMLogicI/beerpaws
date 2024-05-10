package bot

import (
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"strconv"
	"strings"
)

func (b *Bot) setAdditionalPointsHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	values := strings.Split(m.Content, " ")
	if len(values) < 3 {
		_, _ = s.ChannelMessageSend(m.ChannelID, "Введены не все параметры для этой команды!")
		return
	}

	count, err := strconv.Atoi(values[1])
	if err != nil {
		_, _ = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Введено неверное кол-во баллов! : %v", err))
		return
	}

	err = b.setAdditionalPoints(m.Author.ID, values[2], int64(count))
	if err != nil {
		_, _ = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Что-то пошло не так! : %v", err))
		return
	}

	_, _ = s.ChannelMessageSend(m.ChannelID, "Дополнительные очки начислены.")
}

func (b *Bot) setAdditionalPoints(
	discordID string,
	targetDiscordID string,
	count int64,
) error {
	if !b.isAdmin(discordID) {
		return errors.New("вы не можете использовать эту команду")
	}

	targetUser, err := b.userService.GetUserByDiscordID(targetDiscordID)
	if err != nil {
		return err
	}

	if targetUser == nil {
		return errors.New("такой пользователь не зарегистирован")
	}

	return b.pointService.SetAdditionalPoints(targetUser.ID, count)
}
