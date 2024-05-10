package bot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"strings"
)

func (b *Bot) getRulesHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	rules, err := b.pointService.GetPointsRules()
	if err != nil {
		log.Println(err)
		return
	}

	message := strings.Builder{}
	for i, rule := range rules {
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
