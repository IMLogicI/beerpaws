package bot

import (
	"beerpaws/bot/consts"
	"beerpaws/config"
	"beerpaws/service"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

const (
	ruleChunkSize = 10
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
		case strings.HasPrefix(m.Content, consts.GetRulePrefix):
			rules, err := getRules(pointService)
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
		case strings.HasPrefix(m.Content, consts.MakeRequestPrefix):
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
		case strings.HasPrefix(m.Content, consts.AddRulePrefix):
			values := strings.Split(m.Content, " ; ")
			if len(values) < 4 {
				_, _ = s.ChannelMessageSend(m.ChannelID, "Введены не все параметры для этой команды!")
				return
			}

			count, err := strconv.Atoi(values[1])
			if err != nil {
				_, _ = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Введена неверная сумма баллов! : %v", err))
				return
			}

			err = makeNewRule(pointService, userService, m.Author.ID, int64(count), values[2], values[3])
			if err != nil {
				_, _ = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Что-то пошло не так! : %v", err))
				return
			}

			_, _ = s.ChannelMessageSend(m.ChannelID, "Новое правило начисления очков добавлено!")
		case strings.HasPrefix(m.Content, consts.ViewOpenRequestsPrefix):
			requests, err := getOpenedRequests(pointService, userService, m.Author.ID)
			if err != nil {
				_, _ = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Что-то пошло не так! : %v", err))
				return
			}

			for _, request := range requests {
				_, _ = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Номер запроса: %d\n Создатель: @%s\n За что: %s\n Ссылка на скрин: %s\n Подтверждено: %v",
					request.ID,
					request.UserName,
					request.RuleName,
					request.ScreenshotLink,
					request.Approved))
			}
		case strings.HasPrefix(m.Content, consts.ApproveRequestPrefix):
			values := strings.Split(m.Content, " ")
			if len(values) < 2 {
				_, _ = s.ChannelMessageSend(m.ChannelID, "Введены не все параметры для этой команды!")
				return
			}

			requestNumber, err := strconv.Atoi(values[1])
			if err != nil {
				_, _ = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Введен неверный номер запроса! : %v", err))
				return
			}

			err = approveRequest(pointService, userService, m.Author.ID, int64(requestNumber))
			if err != nil {
				_, _ = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Что-то пошло не так! : %v", err))
				return
			}

			_, _ = s.ChannelMessageSend(m.ChannelID, "Запрос подтвержден! Чтобы очки начислились, закройте его.")
		case strings.HasPrefix(m.Content, consts.CloseRequestPrefix):
			values := strings.Split(m.Content, " ")
			if len(values) < 2 {
				_, _ = s.ChannelMessageSend(m.ChannelID, "Введены не все параметры для этой команды!")
				return
			}

			requestNumber, err := strconv.Atoi(values[1])
			if err != nil {
				_, _ = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Введен неверный номер запроса! : %v", err))
				return
			}

			err = closeRequest(pointService, userService, m.Author.ID, int64(requestNumber))
			if err != nil {
				_, _ = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Что-то пошло не так! : %v", err))
				return
			}

			_, _ = s.ChannelMessageSend(m.ChannelID, "Запрос закрыт.")
		case strings.HasPrefix(m.Content, consts.GetMyPointsPrefix):
			count, err := getPointsByDiscordID(m.Author.ID, pointService, userService)
			if err != nil {
				_, _ = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Что-то пошло не так! : %v", err))
				return
			}

			_, _ = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("У вас на баллансе %d баллов", count))
		case strings.HasPrefix(m.Content, consts.DeleteRulePrefix):
			values := strings.Split(m.Content, " ")
			if len(values) < 2 {
				_, _ = s.ChannelMessageSend(m.ChannelID, "Введены не все параметры для этой команды!")
				return
			}

			ruleID, err := strconv.Atoi(values[1])
			if err != nil {
				_, _ = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Введен неверный номер правила! : %v", err))
				return
			}

			err = deleteRule(pointService, userService, m.Author.ID, int64(ruleID))
			if err != nil {
				_, _ = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Что-то пошло не так! : %v", err))
				return
			}

			_, _ = s.ChannelMessageSend(m.ChannelID, "Правило удалено.")
		case strings.HasPrefix(m.Content, consts.HelpPrefix):
			_, _ = s.ChannelMessageSend(m.ChannelID, strings.Join(GetHelpMessage(m.Author.ID, userService), "\n"))
		}
	}
}
