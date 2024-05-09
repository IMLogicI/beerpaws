package bot

import (
	"beerpaws/service"
	"errors"
)

func makePointsRequest(
	discordID string,
	ruleID int64,
	screenshotLink string,
	discordUserName string,
	pointService *service.PointsService,
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
			return errors.New("user not found")
		}
	}

	return pointService.MakePointRequest(user, ruleID, screenshotLink)
}
