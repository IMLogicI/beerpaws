package bot

import (
	"beerpaws/service"
	"errors"
)

func register(
	discordID string,
	discordUserName string,
	userService *service.UserService,
) error {
	user, err := userService.GetUserByDiscordID(discordID)
	if err != nil {
		return err
	}

	if user == nil {
		err = userService.SaveUserFromDiscord(discordID, discordUserName)
		if err != nil {
			return err
		}

		user, err = userService.GetUserByDiscordID(discordID)
		if err != nil {
			return err
		}

		if user == nil {
			return errors.New("user not registered")
		}
	}

	return nil
}
