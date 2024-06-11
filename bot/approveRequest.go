package bot

import (
	"errors"
)

func (b *Bot) approveRequest(
	discordID string,
	requestNumber int64,
) error {
	if !b.isAdmin(discordID) {
		return errors.New("вы не можете использовать эту команду")
	}

	return b.pointService.ApproveRequest(requestNumber)
}
