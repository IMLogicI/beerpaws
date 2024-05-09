package storage

import (
	"beerpaws/storage/consts"
	"beerpaws/storage/models"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
)

type IPoints interface {
	GetPointsRules() ([]models.PointRule, error)
	MakePointRequest(pointRequest models.PointRequest) error
	GetPointsRuleByID(ruleID int64) error
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

func (pointsStorage *PointsStorage) MakePointRequest(pointRequest models.PointRequest) error {
	_, err := pointsStorage.dbConn.Queryx(fmt.Sprintf("%s (%d,%d,'%s')", consts.MakeRequest, pointRequest.RuleID, pointRequest.UserID, pointRequest.ScreenshotLink))
	if err != nil {
		return fmt.Errorf("make point request: %w", err)
	}

	return nil
}

func (pointsStorage *PointsStorage) GetPointsRuleByID(ruleID int64) error {
	log.Printf("%s WHERE '%s'=%d", consts.GetRules, consts.RuleIDField, ruleID)
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
