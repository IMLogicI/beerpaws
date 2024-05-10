package service

import (
	"beerpaws/storage"
	"beerpaws/storage/models"
)

type IUserService interface {
	GetUserByDiscordID(discordID string) (*models.User, error)
	SaveUserFromDiscord(discordID string, discordName string) error
}

type UserService struct {
	userStorage storage.IUser
}

func NewUserService(userStorage *storage.UserStorage) *UserService {
	return &UserService{
		userStorage: userStorage,
	}
}

func (userService *UserService) GetUserByDiscordID(discordID string) (*models.User, error) {
	return userService.userStorage.GetUserByDiscordID(discordID)
}

func (userService *UserService) SaveUserFromDiscord(discordID string, discordName string) error {
	return userService.userStorage.SaveUserFromDiscord(discordID, discordName)
}
