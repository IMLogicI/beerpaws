package bot

import (
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"strconv"
	"strings"
)

func (b *Bot) makePointsRequestHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	values := strings.Split(m.Content, " ")
	if len(values) < 3 {
		_, _ = s.ChannelMessageSend(m.ChannelID, "Введены не все параметры для этой команды!")
		return
	}

	ruleID, err := strconv.Atoi(values[1])
	if err != nil {
		_, _ = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Введен неверный номер правила! : %v", err))
		return
	}

	err = b.makePointsRequest(m.Author.ID, int64(ruleID), values[2], m.Author.Username)
	if err != nil {
		_, _ = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Что-то пошло не так! : %v", err))
		return
	} else {
		_, _ = s.ChannelMessageSend(m.ChannelID, "Запрос отправлен!")
	}
}

func (b *Bot) makePointsRequest(
	discordID string,
	ruleID int64,
	screenshotLink string,
	discordUserName string,
) error {
	user, err := b.userService.GetUserByDiscordID(discordID)
	if err != nil {
		return err
	}

	if user == nil {
		err = b.userService.SaveUserFromDiscord(discordID, discordUserName)
		if err != nil {
			return err
		}

		user, err = b.userService.GetUserByDiscordID(discordID)
		if err != nil {
			return err
		}

		if user == nil {
			return errors.New("user not found")
		}
	}

	return b.pointService.MakePointRequest(user, ruleID, screenshotLink)
}
