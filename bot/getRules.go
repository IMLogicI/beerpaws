package bot

import (
	"beerpaws/storage/models"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"strings"
)

func getRules(s *discordgo.Session, i *discordgo.InteractionCreate, b *Bot) {
	rules, err := b.pointService.GetPointsRules()
	if err != nil {
		log.Println(err)
		return
	}

	earnRules := make([]models.PointRule, 0)
	for _, rule := range rules {
		if rule.Count > 0 {
			earnRules = append(earnRules, rule)
		}
	}

	rulesInteraction(s, i, earnRules)
}

func rulesInteraction(s *discordgo.Session, i *discordgo.InteractionCreate, earnRules []models.PointRule) {
	text := "правила"
	if len(earnRules) > 0 && earnRules[0].Count < 0 {
		text = "лота"
	}
	message := strings.Builder{}

	for _, rule := range earnRules {
		message.WriteString(fmt.Sprintf("Номер %s : %d . %s (%s). %d очков\n", text, rule.ID, rule.Name, rule.Description, rule.Count))
	}

	if message.Len() > 0 {
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags:   discordgo.MessageFlagsEphemeral,
				Content: message.String(),
				Title:   text,
			},
		})
	}
}
