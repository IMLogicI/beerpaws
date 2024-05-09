package service

import (
	"beerpaws/storage"
	"beerpaws/storage/models"
	"errors"
	"time"
)

type IPointsService interface {
	GetPointsRules() ([]models.PointRule, error)
	MakePointRequest(user *models.User, ruleID int64, screenLink string) error
	AddNewRule(newRule models.PointRule) error
	GetOpenedRequests() ([]models.PointRequestForUser, error)
	ApproveRequest(requestID int64) error
	CloseRequest(requestID int64) error
	GetPointsByUserID(userID int64) (int64, error)
	DeleteRule(ruleID int64) error
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

func (pointsService *PointsService) AddNewRule(newRule models.PointRule) error {
	newRule.IsEarned = true
	return pointsService.pointsStorage.AddNewRule(newRule)
}

func (pointsService *PointsService) GetOpenedRequests() ([]models.PointRequestForUser, error) {
	return pointsService.pointsStorage.GetOpenedRequests()
}

func (pointsService *PointsService) ApproveRequest(requestID int64) error {
	return pointsService.pointsStorage.ApproveRequest(requestID)
}

func (pointsService *PointsService) CloseRequest(requestID int64) error {
	request, err := pointsService.pointsStorage.GetRequestByID(requestID)
	if err != nil {
		return err
	}

	if request == nil {
		return errors.New("такого запроса не существует")
	}

	if request.Closed {
		return errors.New("request is already closed")
	}

	if request.Approved {
		err = pointsService.pointsStorage.AddPoints(models.PointHistory{
			RequestID: requestID,
			Time:      time.Now(),
		})

		if err != nil {
			return err
		}
	}

	return pointsService.pointsStorage.CloseRequest(requestID)
}

func (pointsService *PointsService) GetPointsByUserID(userID int64) (int64, error) {
	return pointsService.pointsStorage.GetPointsByUserID(userID)
}

func (pointsService *PointsService) DeleteRule(ruleID int64) error {
	return pointsService.pointsStorage.DeleteRule(ruleID)
}
