package bot

import (
	"beerpaws/bot/consts"
	"beerpaws/config"
	"beerpaws/service"
	"github.com/bwmarrin/discordgo"
	"log"
	"strings"
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
	goBot.AddHandler(b.interactionHandler)

	_, err = goBot.ApplicationCommandCreate(config.ApplicationID, config.GuildID, &discordgo.ApplicationCommand{
		Name:        consts.ButtonInteraction,
		Description: "Добавить кнопку по которой можно отправить запрос",
	})

	_, err = goBot.ApplicationCommandCreate(config.ApplicationID, config.GuildID, &discordgo.ApplicationCommand{
		Name:        consts.GetMyPointsInteraction,
		Description: "Посмотреть мои баллы",
	})

	if err != nil {
		log.Fatalf("Cannot create slash command: %v", err)
	}

	err = goBot.Open()
	if err != nil {
		return
	}
}

func (b *Bot) interactionHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		if h, ok := commandsHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i, b)
		}
	case discordgo.InteractionMessageComponent:
		preffix := i.MessageComponentData().CustomID
		if strings.Contains(i.MessageComponentData().CustomID, consts.AcceptRequestInteraction) {
			preffix = consts.AcceptRequestInteraction
		}

		if strings.Contains(i.MessageComponentData().CustomID, consts.DeclineRequestInteraction) {
			preffix = consts.DeclineRequestInteraction
		}

		if h, ok := componentsHandlers[preffix]; ok {
			h(s, i, b)
		}

	case discordgo.InteractionModalSubmit:
		data := i.ModalSubmitData()

		switch {
		case strings.HasPrefix(data.CustomID, consts.CreateRequestInteraction):
			sendResponsesToChannel(s, i, b)
		case strings.HasPrefix(data.CustomID, consts.CreateSpendRequestInteraction):
			sendSpendResponseToChannel(s, i, b)
		}

	}
}

func (b *Bot) messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if _, ok := config.ChannelsID[m.ChannelID]; ok && m.Author.ID != b.BotID {
		switch {
		case strings.HasPrefix(strings.ToLower(m.Content), consts.AddRulePrefix):
			b.makeNewRuleHandler(s, m)
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

var (
	componentsHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate, b *Bot){
		consts.CreateRequestInteraction:  sendPointRequestForm,
		consts.AcceptRequestInteraction:  acceptRequest,
		consts.DeclineRequestInteraction: declineRequest,
		consts.EarnRulesInteraction:      getRules,
		consts.SpendRulesInteraction:     getSpendRulesHandler,
		consts.GetMyPointsInteraction:    getMyPointsInteraction,
		consts.SpendInteraction:          sendPointSpendForm,
	}
	commandsHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate, b *Bot){
		consts.ButtonInteraction: sendPointRequestButton,
	}
)
