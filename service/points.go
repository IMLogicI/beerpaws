package service

import (
	"beerpaws/domain"
	"beerpaws/storage"
	"beerpaws/storage/models"
	"errors"
	"time"
)

type IPointsService interface {
	GetPointsRules() ([]models.PointRule, error)
	MakePointRequest(user *models.User, ruleID int64, pointsCount int64, screenLink string) (int64, error)
	AddNewRule(newRule domain.PointRule) error
	GetOpenedRequests() ([]models.PointRequestForUser, error)
	ApproveRequest(requestID int64) error
	CloseRequest(requestID int64) error
	GetPointsByUserID(userID int64) (int64, error)
	DeleteRule(ruleID int64) error
	SetAdditionalPoints(userID int64, count int64, reason string) error
	GetRuleByID(ruleID int64) (models.PointRule, error)
}

type PointsService struct {
	pointsStorage storage.IPoints
}

func NewPointsService(pointsStorage *storage.PointsStorage) *PointsService {
	return &PointsService{
		pointsStorage: pointsStorage,
	}
}

func (pointsService *PointsService) GetPointsRules() ([]models.PointRule, error) {
	pointsRules, err := pointsService.pointsStorage.GetPointsRules()
	if err != nil {
		return nil, err
	}

	return pointsRules, nil
}

func (pointsService *PointsService) GetRuleByID(ruleID int64) (models.PointRule, error) {
	return pointsService.pointsStorage.GetPointsRuleByID(ruleID)
}

func (pointsService *PointsService) MakePointRequest(user *models.User, ruleID int64, pointsCount int64, screenLink string) (int64, error) {
	_, err := pointsService.pointsStorage.GetPointsRuleByID(ruleID)
	if err != nil {
		return 0, err
	}

	return pointsService.pointsStorage.MakePointRequest(models.PointRequest{
		RuleID:         ruleID,
		UserID:         user.ID,
		ScreenshotLink: screenLink,
		PointsCount:    pointsCount,
		Approved:       false,
		Closed:         false,
	})
}

func (pointsService *PointsService) AddNewRule(newRule domain.PointRule) error {
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

func (pointsService *PointsService) SetAdditionalPoints(userID int64, count int64, reason string) error {
	id, err := pointsService.pointsStorage.MakePointRequest(models.PointRequest{
		RuleID:         0,
		UserID:         userID,
		ScreenshotLink: reason,
		PointsCount:    count,
		Approved:       true,
		Closed:         true,
	})

	if err != nil {
		return err
	}

	return pointsService.pointsStorage.AddPoints(models.PointHistory{
		RequestID: id,
		Time:      time.Now(),
	})
}
