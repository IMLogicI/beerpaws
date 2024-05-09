package service

import (
	"beerpaws/storage"
	"beerpaws/storage/models"
)

type IPointsService interface {
	GetPointsRules() ([]models.PointRule, error)
	MakePointRequest(user *models.User, ruleID int64, screenLink string) error
}

type PointsService struct {
	pointsStorage *storage.PointsStorage
}

func NewPointsService(pointsStorage *storage.PointsStorage) *PointsService {
	return &PointsService{
		pointsStorage: pointsStorage,
	}
}

func (pointsService *PointsService) GetPointsRules() ([]models.PointRule, error) {
	return pointsService.pointsStorage.GetPointsRules()
}

func (pointsService *PointsService) MakePointRequest(user *models.User, ruleID int64, screenLink string) error {
	err := pointsService.pointsStorage.GetPointsRuleByID(ruleID)
	if err != nil {
		return err
	}

	return pointsService.pointsStorage.MakePointRequest(models.PointRequest{
		RuleID:         ruleID,
		UserID:         user.ID,
		ScreenshotLink: screenLink,
	})
}
