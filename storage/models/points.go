package models

import (
	"database/sql"
	"time"
)

type PointRule struct {
	ID          int64         `db:"id"`
	Name        string        `db:"name"`
	Description string        `db:"description"`
	Count       int64         `db:"count"`
	IsEarned    bool          `db:"is_earned"`
	DaysActual  sql.NullInt64 `db:"days_actual"`
}

type PointRequest struct {
	ID             int64  `db:"id"`
	RuleID         int64  `db:"rule_id"`
	UserID         int64  `db:"user_id"`
	ScreenshotLink string `db:"screenshot_link"`
	Approved       bool   `db:"approved"`
	Closed         bool   `db:"closed"`
	PointsCount    int64  `db:"points_count"`
}

type PointRequestForUser struct {
	ID             int64  `db:"id"`
	RuleName       string `db:"name"`
	UserName       string `db:"nickname"`
	ScreenshotLink string `db:"screenshot_link"`
	Approved       bool   `db:"approved"`
	PointsCount    int64  `db:"points_count"`
}

type PointHistory struct {
	ID        int64     `db:"id"`
	RequestID int64     `db:"request_id"`
	Time      time.Time `db:"time"`
}
