package bot

import (
	"beerpaws/bot/consts"
	"beerpaws/config"
	"beerpaws/service"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
)

const (
	ruleChunkSize = 10
)

type Bot struct {
	BotID        string
	pointService service.IPointsService
	userService  service.IUserService
}

func NewBot(pService service.IPointsService, uService service.IUserService) *Bot {
	return &Bot{
		pointService: pService,
		userService:  uService,
	}
}

func (b *Bot) Run() {
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

	b.BotID = user.ID
	goBot.AddHandler(b.messageHandler)
	err = goBot.Open()
	if err != nil {
		return
	}
}

func (b *Bot) messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if _, ok := config.ChannelsID[m.ChannelID]; ok && m.Author.ID != b.BotID {
		switch {
		case strings.HasPrefix(strings.ToLower(m.Content), consts.GetRulePrefix):
			b.getRulesHandler(s, m)
		case strings.HasPrefix(strings.ToLower(m.Content), consts.MakeRequestPrefix):
			b.makePointsRequestHandler(s, m)
		case strings.HasPrefix(strings.ToLower(m.Content), consts.AddRulePrefix):
			b.makeNewRuleHandler(s, m)
		case strings.HasPrefix(strings.ToLower(m.Content), consts.ViewOpenRequestsPrefix):
			b.getOpenedRequestsHandler(s, m)
		case strings.HasPrefix(strings.ToLower(m.Content), consts.ApproveRequestPrefix):
			b.approveRequestHandler(s, m)
		case strings.HasPrefix(strings.ToLower(m.Content), consts.CloseRequestPrefix):
			b.closeRequestHandler(s, m)
		case strings.HasPrefix(strings.ToLower(m.Content), consts.GetMyPointsPrefix):
			b.getPointsByDiscordIDHandler(s, m)
		case strings.HasPrefix(strings.ToLower(m.Content), consts.DeleteRulePrefix):
			b.deleteRuleHandler(s, m)
		case strings.HasPrefix(strings.ToLower(m.Content), consts.HelpPrefix):
			b.GetHelpMessageHandler(s, m)
		case strings.HasPrefix(strings.ToLower(m.Content), consts.SetAdditionalPointPrefix):
			b.setAdditionalPointsHandler(s, m)
		case strings.HasPrefix(strings.ToLower(m.Content), consts.RegisterPrefix):
			b.registerHandler(s, m)
		}
	}
}
