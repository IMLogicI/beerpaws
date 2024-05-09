package bot

import (
	"beerpaws/service"
	"beerpaws/storage/models"
	"errors"
)

func getOpenedRequests(
	pointService *service.PointsService,
	userService *service.UserService,
	discordID string,
) ([]models.PointRequestForUser, error) {
	if !isAdmin(userService, discordID) {
		return nil, errors.New("вы не можете использовать эту команду")
	}

	return pointService.GetOpenedRequests()
}
