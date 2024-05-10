package bot

import (
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
)

func (b *Bot) registerHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	err := b.register(m.Author.ID, m.Author.Username)
	if err != nil {
		_, _ = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Что-то пошло не так! : %v", err))
		return
	}
	_, _ = s.ChannelMessageSend(m.ChannelID, "Вы зарегистрированы в системе")
}

func (b *Bot) register(
	discordID string,
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
			return errors.New("user not registered")
		}
	}

	return nil
}
