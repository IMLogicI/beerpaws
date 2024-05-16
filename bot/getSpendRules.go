package bot

import (
	"beerpaws/storage/models"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"strings"
)

func (b *Bot) getSpendRulesHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
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

	message := strings.Builder{}
	for i, rule := range spendRules {
		message.WriteString(fmt.Sprintf("Номер правила : %d . %s (%s). %d очков\n", rule.ID, rule.Name, rule.Description, rule.Count))
		if (i+1)%ruleChunkSize == 0 {
			_, _ = s.ChannelMessageSend(m.ChannelID, message.String())
			message = strings.Builder{}
		}
	}

	if message.Len() > 0 {
		_, _ = s.ChannelMessageSend(m.ChannelID, message.String())
	}
}
