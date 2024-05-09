package bot

import "beerpaws/service"

func getPointsByDiscordID(
	discordID string,
	pointService *service.PointsService,
	userService *service.UserService,
) (int64, error) {
	user, err := userService.GetUserByDiscordID(discordID)
	if err != nil {
		return 0, err
	}

	return pointService.GetPointsByUserID(user.ID)
}
