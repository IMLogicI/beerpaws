package bot

import (
	"beerpaws/service"
	"errors"
)

func setAdditionalPoints(
	pointService *service.PointsService,
	userService *service.UserService,
	discordID string,
	targetDiscordID string,
	count int64,
) error {
	if !isAdmin(userService, discordID) {
		return errors.New("вы не можете использовать эту команду")
	}

	targetUser, err := userService.GetUserByDiscordID(targetDiscordID)
	if err != nil {
		return err
	}

	if targetUser == nil {
		return errors.New("такой пользователь не зарегистирован")
	}

	return pointService.SetAdditionalPoints(targetUser.ID, count)
}
