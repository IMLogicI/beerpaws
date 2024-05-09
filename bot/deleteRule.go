package bot

import (
	"beerpaws/service"
	"errors"
)

func deleteRule(
	pointService *service.PointsService,
	userService *service.UserService,
	discordID string,
	ruleID int64,
) error {
	if !isAdmin(userService, discordID) {
		return errors.New("вы не можете использовать эту команду")
	}

	return pointService.DeleteRule(ruleID)
}
