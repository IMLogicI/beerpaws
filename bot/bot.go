package bot

import (
	"beerpaws/bot/consts"
	"beerpaws/config"
	"beerpaws/service"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var (
	BotID        string
	pointService *service.PointsService
	userService  *service.UserService
)

func Run(pService *service.PointsService, uService *service.UserService) {
	pointService = pService
	userService = uService

	// create bot session
	goBot, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		log.Fatal(err)
		return
	}
	// make the bot a user
	user, err := goBot.User("@me")
	if err != nil {
		log.Fatal(err)
		return
	}

	BotID = user.ID
	goBot.AddHandler(messageHandler)
	err = goBot.Open()
	if err != nil {
		return
	}
}

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID != BotID {
		switch {
		case strings.HasPrefix(m.Content, consts.GetRulePreffix):
			rules, err := pointService.GetPointsRules()
			if err != nil {
				log.Println(err)
				return
			}

			j, _ := json.Marshal(rules)

			_, _ = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s", string(j)))
		case strings.HasPrefix(m.Content, consts.MakeRequestPreffix):
			values := strings.Split(m.Content, " ")
			if len(values) < 3 {
				_, _ = s.ChannelMessageSend(m.ChannelID, "Введены не все параметры для этой команды!")
				return
			}

			ruleID, err := strconv.Atoi(values[1])
			if err != nil {
				_, _ = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Введен неверный номер правила! : %v", err))
				return
			}

			err = makePointsRequest(m.Author.ID, int64(ruleID), values[2], m.Author.Username, pointService, userService)
			if err != nil {
				_, _ = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Что-то пошло не так! : %v", err))
				return
			} else {
				_, _ = s.ChannelMessageSend(m.ChannelID, "Запрос отправлен!")
			}
		}
	}
}
