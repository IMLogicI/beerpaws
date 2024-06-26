package bot

import (
	"github.com/bwmarrin/discordgo"
	"strings"
)

var (
	userHelp = []string{
		"**!getPointRules** - Посмотреть вещи, за которые можно получить баллы",
		"**!getSpendRules** - Посмотреть вещи, на которые можно потратить баллы",
		"**!makeEarnRequest** <Номер правила> <Сколько очков должно быть начислено> <Ссылка на скриншот> - Подать заявку на выдачу или трату очков (через пробел укажи номер правила, сколько очков должно быть начислено и ссылку на скриншот), ",
		"**!getMyPoints** - Посмотреть количество очков на балансе",
		"**!register** - Зарегистрироваться в системе (понадобится для занесения начальных очков, далее авторегистрация после первого запроса)",
	}
	adminHelp = []string{
		"**!addEarnRule** ; <Количество начисляемых очков> ; <Название правила> ; <Описание правила> ; <Колво дней,сколько действует> - Создать новое правило (кол-во дней для случая, когда правило на трату должно неделю висеть на игроке)",
		"**!viewOpenRequests** - Посмотреть все открытые запросы, ожидающие рассмотрения",
		"**!approveRequest** <Номер запроса> - Поддтвердить, что завпрос верен",
		"**!closeRequest** <Номер запроса> - Закрыть подтверженный запрос с начислением очков, либо неподтвержденный без начисления",
		"**!deleteRule** <Номер правила> - Удалить правило",
		"**!setAdditionalPoint** ; <Кол-во очков> ;  <Discord ID> ; <За что> - Начислить дополнительные очки пользователю",
	}
)

func (b *Bot) GetHelpMessageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	helps := userHelp
	if b.isAdmin(m.Author.ID) {
		helps = append(helps, adminHelp...)
	}

	_, _ = s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
		Title:       "Список команд",
		Description: strings.Join(helps, "\n"),
	})
}
