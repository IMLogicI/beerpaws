package bot

import (
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"strconv"
	"strings"
)

func (b *Bot) approveRequestHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	values := strings.Split(m.Content, " ")
	if len(values) < 2 {
		_, _ = s.ChannelMessageSend(m.ChannelID, "Введены не все параметры для этой команды!")
		return
	}

	requestNumber, err := strconv.Atoi(values[1])
	if err != nil {
		_, _ = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Введен неверный номер запроса! : %v", err))
		return
	}

	err = b.approveRequest(m.Author.ID, int64(requestNumber))
	if err != nil {
		_, _ = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Что-то пошло не так! : %v", err))
		return
	}

	_, _ = s.ChannelMessageSend(m.ChannelID, "Запрос подтвержден! Чтобы очки начислились, закройте его.")
}

func (b *Bot) approveRequest(
	discordID string,
	requestNumber int64,
) error {
	if !b.isAdmin(discordID) {
		return errors.New("вы не можете использовать эту команду")
	}

	return b.pointService.ApproveRequest(requestNumber)
}
