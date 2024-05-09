package bot

import (
	"beerpaws/service"
	"errors"
)

func closeRequest(
	pointService *service.PointsService,
	userService *service.UserService,
	discordID string,
	requestNumber int64,
) error {
	if !isAdmin(userService, discordID) {
		return errors.New("вы не можете использовать эту команду")
	}

	return pointService.CloseRequest(requestNumber)
}
