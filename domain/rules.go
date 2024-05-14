package domain

type PointRule struct {
	ID          int64
	Name        string
	Description string
	Count       int64
	IsEarned    bool
	DaysActual  *int64
}
