package models

type PointRule struct {
	ID          int64  `db:"id"`
	Name        string `db:"name"`
	Description string `db:"description"`
	Count       int64  `db:"count"`
	IsEarned    bool   `db:"is_earned"`
}

type PointRequest struct {
	ID             int64  `db:"id"`
	RuleID         int64  `db:"rule_id"`
	UserID         int64  `db:"user_id"`
	ScreenshotLink string `db:"screenshot_link"`
	Approved       bool   `db:"approved"`
	Closed         bool   `db:"closed"`
}
