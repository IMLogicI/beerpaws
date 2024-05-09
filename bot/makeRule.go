package bot

import (
	"beerpaws/domain"
	"beerpaws/service"
	"beerpaws/storage/models"
	"errors"
	"log"
)

func makeNewRule(
	pointService *service.PointsService,
	userService *service.UserService,
	discordID string,
	count int64,
	name string,
	description string,
) error {
	if !isAdmin(userService, discordID) {
		return errors.New("вы не можете использовать эту команду")
	}

	return pointService.AddNewRule(models.PointRule{
		Name:        name,
		Description: description,
		Count:       count,
	})
}

func isAdmin(userService *service.UserService, discordID string) bool {
	user, err := userService.GetUserByDiscordID(discordID)
	if err != nil {
		log.Printf(err.Error())
		return false
	}

	if user == nil || user.Role != domain.AccessAdmin {
		return false
	}

	return true
}
