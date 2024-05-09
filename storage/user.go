package storage

import (
	"beerpaws/domain"
	"beerpaws/storage/consts"
	"beerpaws/storage/models"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type IUser interface {
	GetUserByDiscordID(discordID int64) (*models.User, error)
	SaveUserFromDiscord(discordID int64, discordName string) error
}

type UserStorage struct {
	dbConn *sqlx.DB
}

func NewUserStorage(dbConn *sqlx.DB) *UserStorage {
	return &UserStorage{dbConn: dbConn}
}

func (usersStorage *UserStorage) GetUserByDiscordID(discordID string) (*models.User, error) {
	rows, err := usersStorage.dbConn.Queryx(fmt.Sprintf("%s'%s'", consts.GetUserByDiscordID, discordID))
	if err != nil {
		return nil, fmt.Errorf("get user from database: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var user models.User
		if err := rows.StructScan(&user); err != nil {
			return nil, fmt.Errorf("get user from database: %w", err)
		}

		return &user, nil
	}

	return nil, nil
}

func (usersStorage *UserStorage) SaveUserFromDiscord(discordID string, discordName string) error {
	_, err := usersStorage.dbConn.Queryx(fmt.Sprintf("%s ('%s','%s','%s', '%s')", consts.SaveUserFromDiscord, discordName, discordID, discordName, domain.AccessUser))
	if err != nil {
		return fmt.Errorf("save user from discord: %w", err)
	}

	return nil
}
