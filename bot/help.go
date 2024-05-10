package bot

import (
	"github.com/bwmarrin/discordgo"
	"strings"
)

var (
	userHelp = []string{
		"!getPointRules - Посмотреть правила, за которые можно получить баллы",
		"!makeEarnRequest <Номер правила> <Ссылка на скриншот> - Подать заявку на выдачу очков (через пробел укажи номер правила и ссылку на скриншот)",
		"!getMyPoints - Посмотреть количество очков на баллансе",
		"!register - Зарегистрироваться в системе (понадобится для занесения начальных очков, далее авторегистрация после первого запроса)",
	}
	adminHelp = []string{"!addEarnRule ; <Еоличество начисляемых очков> ; <Название правила> ; <Описание правила> - Создать новое правило",
		"!viewOpenRequests - Посмотреть все открытые запросы, ожидающие рассмотрения",
		"!approveRequest <Номер запроса> - Поддтвердить, что завпрос верен",
		"!closeRequest <Номер запроса> - Закрыть подтверженный запрос с начислением очков, либо неподтвержденный без начисления",
		"!deleteRule <Номер правила> - Удалить правило",
		"!setAdditionalPoint <Кол-во очков> <Discord ID> - Начислить дополнительные очки пользователю",
	}
)

func (b *Bot) GetHelpMessageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	helps := userHelp
	if b.isAdmin(m.Author.ID) {
		helps = append(helps, adminHelp...)
	}

	_, _ = s.ChannelMessageSend(m.ChannelID, strings.Join(helps, "\n"))
}
