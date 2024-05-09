package bot

import (
	"beerpaws/service"
	"beerpaws/storage/models"
)

func getRules(pointService *service.PointsService) ([]models.PointRule, error) {
	return pointService.GetPointsRules()
}
