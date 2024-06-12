package bot

import (
	"beerpaws/domain"
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"strconv"
	"strings"
)

func (b *Bot) makeNewRuleHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	values := strings.Split(m.Content, " ; ")
	if len(values) < 4 {
		_, _ = s.ChannelMessageSend(m.ChannelID, "Введены не все параметры для этой команды!")
		return
	}

	count, err := strconv.Atoi(values[1])
	if err != nil {
		_, _ = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Введена неверная сумма баллов! : %v", err))
		return
	}

	err = b.makeNewRule(m.Author.ID, int64(count), values[2], values[3])
	if err != nil {
		_, _ = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Что-то пошло не так! : %v", err))
		return
	}

	_, _ = s.ChannelMessageSend(m.ChannelID, "Новое правило начисления очков добавлено!")
}

func (b *Bot) makeNewRule(
	discordID string,
	count int64,
	name string,
	description string,
) error {
	if !b.isAdmin(discordID) {
		return errors.New("вы не можете использовать эту команду")
	}

	return b.pointService.AddNewRule(domain.PointRule{
		Name:        name,
		Description: description,
		Count:       count,
	})
}

func (b *Bot) isAdmin(discordID string) bool {
	user, err := b.userService.GetUserByDiscordID(discordID)
	if err != nil {
		log.Printf(err.Error())
		return false
	}

	if user == nil || user.Role != domain.AccessAdmin {
		return false
	}

	return true
}
