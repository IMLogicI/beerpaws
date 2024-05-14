package storage

import (
	"beerpaws/domain"
	"beerpaws/storage/consts"
	"beerpaws/storage/models"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
)

const (
	timeLayout = "2006-01-02 15:04:05"
)

type IPoints interface {
	GetPointsRules() ([]models.PointRule, error)
	MakePointRequest(pointRequest models.PointRequest) (int64, error)
	GetPointsRuleByID(ruleID int64) error
	AddNewRule(newRule domain.PointRule) error
	GetOpenedRequests() ([]models.PointRequestForUser, error)
	ApproveRequest(requestID int64) error
	AddPoints(pointsAdding models.PointHistory) error
	GetRequestByID(requestID int64) (*models.PointRequest, error)
	CloseRequest(requestID int64) error
	GetPointsByUserID(userID int64) (int64, error)
	DeleteRule(ruleID int64) error
	GetAdditionalPointsFromUserID(user int64) (int64, error)
	SetAdditionalPoints(userID int64, count int64) error
}

type PointsStorage struct {
	dbConn *sqlx.DB
}

func NewPointsStorage(dbConn *sqlx.DB) *PointsStorage {
	return &PointsStorage{dbConn: dbConn}
}

func (pointsStorage *PointsStorage) GetPointsRules() ([]models.PointRule, error) {
	rows, err := pointsStorage.dbConn.Queryx(consts.GetRules)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	rules := make([]models.PointRule, 0)
	for rows.Next() {
		var rule models.PointRule
		if err := rows.StructScan(&rule); err != nil {
			return nil, err
		}
		rules = append(rules, rule)
	}

	// Check for errors from iterating over rows.
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return rules, nil
}

func (pointsStorage *PointsStorage) MakePointRequest(pointRequest models.PointRequest) (int64, error) {
	rows, err := pointsStorage.dbConn.Queryx(fmt.Sprintf("%s (%d,%d,'%s',%v,%v,%d) RETURNING id", consts.MakeRequest, pointRequest.RuleID, pointRequest.UserID, pointRequest.ScreenshotLink, pointRequest.Approved, pointRequest.Closed, pointRequest.PointsCount))
	if err != nil {
		return 0, fmt.Errorf("make point request: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var requestID int64
		err := rows.Scan(&requestID)
		return requestID, err
	}

	return 0, nil
}

func (pointsStorage *PointsStorage) GetPointsRuleByID(ruleID int64) error {
	rows, err := pointsStorage.dbConn.Queryx(fmt.Sprintf("%s WHERE \"%s\"=%d", consts.GetRules, consts.IDField, ruleID))
	if err != nil {
		return err
	}

	defer rows.Close()
	for rows.Next() {
		return nil
	}

	return errors.New("unknown rule")
}

func (pointsStorage *PointsStorage) AddNewRule(newRule domain.PointRule) error {
	_, err := pointsStorage.dbConn.Queryx(fmt.Sprintf("%s (%d,'%s','%s',%v,%d)", consts.AddNewRule, newRule.Count, newRule.Name, newRule.Description, newRule.IsEarned, *newRule.DaysActual))
	if err != nil {
		return fmt.Errorf("add new rule: %w", err)
	}

	return nil
}

func (pointsStorage *PointsStorage) GetOpenedRequests() ([]models.PointRequestForUser, error) {
	rows, err := pointsStorage.dbConn.Queryx(consts.GetOpenedRequest)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	requests := make([]models.PointRequestForUser, 0)
	for rows.Next() {
		var request models.PointRequestForUser
		if err := rows.StructScan(&request); err != nil {
			return nil, err
		}
		requests = append(requests, request)
	}

	// Check for errors from iterating over rows.
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return requests, nil
}

func (pointsStorage *PointsStorage) ApproveRequest(requestID int64) error {
	_, err := pointsStorage.dbConn.Queryx(fmt.Sprintf("%s%d", consts.ApproveRequest, requestID))
	if err != nil {
		return fmt.Errorf("approve request: %w", err)
	}

	return nil
}

func (pointsStorage *PointsStorage) CloseRequest(requestID int64) error {
	_, err := pointsStorage.dbConn.Queryx(fmt.Sprintf("%s%d", consts.CloseRequest, requestID))
	if err != nil {
		return fmt.Errorf("close request: %w", err)
	}

	return nil
}

func (pointsStorage *PointsStorage) AddPoints(pointsAdding models.PointHistory) error {
	_, err := pointsStorage.dbConn.Queryx(fmt.Sprintf(consts.AddPoints, pointsAdding.RequestID, pointsAdding.Time.Format(timeLayout)))
	if err != nil {
		return fmt.Errorf("add points: %w", err)
	}

	return nil
}

func (pointsStorage *PointsStorage) GetRequestByID(requestID int64) (*models.PointRequest, error) {
	rows, err := pointsStorage.dbConn.Queryx(fmt.Sprintf(consts.GetRequestByID, requestID))
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var request models.PointRequest
		err := rows.StructScan(&request)
		return &request, err
	}

	return nil, nil
}

func (pointsStorage *PointsStorage) GetPointsByUserID(userID int64) (int64, error) {
	rows, err := pointsStorage.dbConn.Queryx(fmt.Sprintf(consts.GetPoints, userID))
	if err != nil {
		return 0, err
	}

	defer rows.Close()
	for rows.Next() {
		var count int64
		err := rows.Scan(&count)
		return count, err
	}

	return 0, nil
}

func (pointsStorage *PointsStorage) DeleteRule(ruleID int64) error {
	_, err := pointsStorage.dbConn.Queryx(fmt.Sprintf(consts.DeleteRule, ruleID))
	if err != nil {
		return fmt.Errorf("delete rule: %w", err)
	}

	return nil
}

func (pointsStorage *PointsStorage) GetAdditionalPointsFromUserID(user int64) (int64, error) {
	rows, err := pointsStorage.dbConn.Queryx(fmt.Sprintf(consts.GetAdditionalPoints, user))
	if err != nil {
		return 0, fmt.Errorf("get additional points: %w", err)
	}

	defer rows.Close()
	for rows.Next() {
		var count int64
		err := rows.Scan(&count)
		return count, err
	}

	return 0, nil
}

func (pointsStorage *PointsStorage) SetAdditionalPoints(userID int64, count int64) error {
	additionalPoints, err := pointsStorage.GetAdditionalPointsFromUserID(userID)
	if err != nil {
		log.Printf("get additional points: %v", err)
	}

	_, err = pointsStorage.dbConn.Queryx(fmt.Sprintf(consts.SetAdditionalPoints, userID, count, count+additionalPoints))
	if err != nil {
		return fmt.Errorf("get additional points: %w", err)
	}

	return nil
}
