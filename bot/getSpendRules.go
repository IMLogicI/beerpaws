package bot

import (
	"beerpaws/storage/models"
	"github.com/bwmarrin/discordgo"
	"log"
)

func getSpendRulesHandler(s *discordgo.Session, i *discordgo.InteractionCreate, b *Bot) {
	rules, err := b.pointService.GetPointsRules()
	if err != nil {
		log.Println(err)
		return
	}

	spendRules := make([]models.PointRule, 0)
	for _, rule := range rules {
		if rule.Count < 0 {
			spendRules = append(spendRules, rule)
		}
	}

	rulesInteraction(s, i, spendRules)
}
