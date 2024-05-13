package bot

import (
	"beerpaws/storage/models"
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
)

func (b *Bot) getOpenedRequestsHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	requests, err := b.getOpenedRequests(m.Author.ID)
	if err != nil {
		_, _ = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Что-то пошло не так! : %v", err))
		return
	}

	for _, request := range requests {
		_, _ = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Номер запроса: %d\n Создатель: @%s\n За что: %s\n Сколько баллов хочет: %d\n Ссылка на скрин: %s\n Подтверждено: %v\n ====================================",
			request.ID,
			request.UserName,
			request.RuleName,
			request.PointsCount,
			request.ScreenshotLink,
			request.Approved))
	}
}

func (b *Bot) getOpenedRequests(
	discordID string,
) ([]models.PointRequestForUser, error) {
	if !b.isAdmin(discordID) {
		return nil, errors.New("вы не можете использовать эту команду")
	}

	return b.pointService.GetOpenedRequests()
}
