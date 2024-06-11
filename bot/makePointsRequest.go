package bot

import (
	"beerpaws/bot/consts"
	"beerpaws/config"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func (b *Bot) makePointsRequest(
	discordID string,
	ruleID int64,
	pointsCount int64,
	screenshotLink string,
	discordUserName string,
) (int64, error) {
	user, err := b.userService.GetUserByDiscordID(discordID)
	if err != nil {
		return 0, err
	}

	if user == nil {
		err = b.userService.SaveUserFromDiscord(discordID, discordUserName)
		if err != nil {
			return 0, err
		}

		user, err = b.userService.GetUserByDiscordID(discordID)
		if err != nil {
			return 0, err
		}

		if user == nil {
			return 0, errors.New("user not found")
		}
	}

	return b.pointService.MakePointRequest(user, ruleID, pointsCount, screenshotLink)
}

func sendPointRequestButton(s *discordgo.Session, i *discordgo.InteractionCreate, _ *Bot) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Чтобы подать запрос на получение баллов нажми на кнопку ниже.",
			// Buttons and other components are specified in Components field.
			Components: []discordgo.MessageComponent{
				// ActionRow is a container of all buttons within the same row.
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							// Label is what the user will see on the button.
							Label: "Подать запрос на получение очков",
							// Style provides coloring of the button. There are not so many styles tho.
							Style: discordgo.SuccessButton,
							// Disabled allows bot to disable some buttons for users.
							Disabled: false,
							// CustomID is a thing telling Discord which data to send when this button will be pressed.
							CustomID: consts.CreateRequestInteraction,
						},
					},
				},
			},
		},
	})
	if err != nil {
		log.Print(err)
	}
}

func sendPointRequestForm(s *discordgo.Session, i *discordgo.InteractionCreate, _ *Bot) {
	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID: fmt.Sprintf("%s_%s", consts.CreateRequestInteraction, i.Interaction.Member.User.ID),
			Title:    "Создать запрос на баллы",
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:    "rule",
							Label:       "Правило, за которое нужно начислить очки:",
							Style:       discordgo.TextInputShort,
							Placeholder: "1",
							Required:    true,
							MaxLength:   30,
						},
					},
				},
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:  "screenshort",
							Label:     "Ссылка на cкриншот",
							Style:     discordgo.TextInputParagraph,
							Required:  true,
							MaxLength: 2000,
						},
					},
				},
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:  "count",
							Label:     "Сколько очков нужно начислить: ",
							Style:     discordgo.TextInputShort,
							Required:  true,
							MaxLength: 30,
						},
					},
				},
			},
		},
	}
	err := s.InteractionRespond(i.Interaction, response)
	if err != nil {
		log.Print(err)
	}
}

func sendResponsesToChannel(s *discordgo.Session, i *discordgo.InteractionCreate, b *Bot) {
	data := i.ModalSubmitData()

	if !strings.HasPrefix(data.CustomID, consts.CreateRequestInteraction) {
		return
	}

	userid := strings.Split(data.CustomID, "_")[2]

	ruleNumber, errRule := strconv.Atoi(data.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value)
	pointsCount, errCount := strconv.Atoi(data.Components[2].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value)
	screenLinks := data.Components[1].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value

	if errRule != nil || errCount != nil {
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags:   discordgo.MessageFlagsEphemeral,
				Content: "Введены неверные данные!",
				Title:   "Введены неверные данные!",
			},
		})

		return
	}

	id, err := b.makePointsRequest(userid, int64(ruleNumber), int64(pointsCount), screenLinks, i.Interaction.Member.User.Username)
	if err != nil {
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags:   discordgo.MessageFlagsEphemeral,
				Content: "Что-то пошло не так!",
				Title:   "Что-то пошло не так!",
			},
		})

		return
	}

	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:   discordgo.MessageFlagsEphemeral,
			Content: "Запрос отправлен!",
			Title:   "Запрос отправлен!",
		},
	})

	i.Interaction.ChannelID = config.ChannelsAdminID
	_, err = s.ChannelMessageSendComplex(config.ChannelsAdminID, &discordgo.MessageSend{
		Content: fmt.Sprintf(
			"Запрос на баллы от <@%s>\n\n**Правило**:\n%v\n\n**Скриншоты**:\n%s\n\n**Сколько очков просит**:\n%v",
			userid,
			ruleNumber,
			screenLinks,
			pointsCount,
		),
		// Buttons and other components are specified in Components field.
		Components: []discordgo.MessageComponent{
			// ActionRow is a container of all buttons within the same row.
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						// Label is what the user will see on the button.
						Label: "Подтвердить",
						// Style provides coloring of the button. There are not so many styles tho.
						Style: discordgo.SuccessButton,
						// Disabled allows bot to disable some buttons for users.
						Disabled: false,
						// CustomID is a thing telling Discord which data to send when this button will be pressed.
						CustomID: fmt.Sprintf("%s_%v", consts.AcceptRequestInteraction, id),
					},
					discordgo.Button{
						// Label is what the user will see on the button.
						Label: "Отклонить",
						// Style provides coloring of the button. There are not so many styles tho.
						Style: discordgo.DangerButton,
						// Disabled allows bot to disable some buttons for users.
						Disabled: false,
						// CustomID is a thing telling Discord which data to send when this button will be pressed.
						CustomID: fmt.Sprintf("%s_%v", consts.DeclineRequestInteraction, id),
					},
				},
			},
		},
	})
	if err != nil {
		log.Print(err)
	}
}

func acceptRequest(s *discordgo.Session, i *discordgo.InteractionCreate, b *Bot) {
	requestID := strings.Split(i.MessageComponentData().CustomID, "_")[2]
	reqIDNum, err := strconv.Atoi(requestID)
	if err != nil {
		log.Print(err)
		return
	}
	userID := i.Member.User.ID

	err = b.approveRequest(userID, int64(reqIDNum))
	if err != nil {
		errInteraction(s, i)
		return
	}

	err = b.closeRequest(userID, int64(reqIDNum))
	if err != nil {
		errInteraction(s, i)
		return
	}

	disableButtons(s, i)
	successInteraction(s, i)
}

func declineRequest(s *discordgo.Session, i *discordgo.InteractionCreate, b *Bot) {
	requestID := strings.Split(i.MessageComponentData().CustomID, "_")[2]
	reqIDNum, err := strconv.Atoi(requestID)
	if err != nil {
		log.Print(err)
		return
	}

	userID := i.Member.User.ID

	err = b.closeRequest(userID, int64(reqIDNum))
	if err != nil {
		errInteraction(s, i)
		return
	}

	disableButtons(s, i)
	successInteraction(s, i)
}

func disableButtons(s *discordgo.Session, i *discordgo.InteractionCreate) {
	comps := i.Message.Components
	comps[0].(*discordgo.ActionsRow).Components[0].(*discordgo.Button).Disabled = true
	comps[0].(*discordgo.ActionsRow).Components[1].(*discordgo.Button).Disabled = true

	_, _ = s.ChannelMessageEditComplex(&discordgo.MessageEdit{
		Content:    &i.Message.Content,
		Components: &comps,
		ID:         i.Message.ID,
		Channel:    i.ChannelID,
	})
}

func errInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:   discordgo.MessageFlagsEphemeral,
			Content: "Что-то пошло не так!",
			Title:   "Что-то пошло не так!",
		},
	})
}

func successInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:   discordgo.MessageFlagsEphemeral,
			Content: "Готово!",
			Title:   "Готово!",
		},
	})
}
