package server

import (
	"errors"
	"github.com/bwmarrin/discordgo"
)

type Bot struct {
	api     *discordgo.Session
	service *Service
	cmds    []*discordgo.ApplicationCommand
}

type CustomCommandHandler func(i *discordgo.InteractionCreate, options map[string]*discordgo.ApplicationCommandInteractionDataOption) string

func NewBot(discordToken string, service *Service) (*Bot, error) {
	api, err := discordgo.New("Bot " + discordToken)
	if err != nil {
		return nil, err
	}
	err = api.Open()
	if err != nil {
		return nil, err
	}
	return &Bot{api: api, service: service}, nil
}

func (b *Bot) RegisterCommands(guildId string) {
	b.cmds = make([]*discordgo.ApplicationCommand, 0)
	cmds := []*discordgo.ApplicationCommand{
		{
			Name:        "bind",
			Description: "Bind your minecraft to your discord account",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "code",
					Description: "in-game minecraft account code",
					Required:    true,
				},
			},
		},
		{
			Name:        "unbind",
			Description: "Unbind your minecraft account from your discord account",
		},
	}
	handlers := map[string]CustomCommandHandler{
		"bind": func(i *discordgo.InteractionCreate, options map[string]*discordgo.ApplicationCommandInteractionDataOption) string {
			codeInfo, err := b.service.CheckCode(options["code"].StringValue())
			if err != nil {
				if errors.As(err, &ApplicationError{}) {
					return err.Error()
				} else {
					return "Something went wrong"
				}
			}
			discordId := i.Member.User.ID
			_, err = b.service.CreateUser(discordId, codeInfo.XUID)
			if err != nil {
				if errors.As(err, &ApplicationError{}) {
					return err.Error()
				} else {
					return "Something went wrong"
				}
			}
			_ = b.service.RevokeCode(codeInfo.Code)
			return "Binding has been created successfully"
		},
		"unbind": func(i *discordgo.InteractionCreate, options map[string]*discordgo.ApplicationCommandInteractionDataOption) string {
			discordId := i.Member.User.ID
			err := b.service.DeleteUserByDiscord(discordId)
			if err != nil {
				if errors.As(err, &ApplicationError{}) {
					return err.Error()
				} else {
					return "Something went wrong"
				}
			}
			return "Binding has been removed successfully"
		},
	}

	for _, cmd := range cmds {
		registeredCmd, err := b.api.ApplicationCommandCreate(b.api.State.User.ID, guildId, cmd)
		if err != nil {
			panic(err)
		}
		b.cmds = append(b.cmds, registeredCmd)
	}

	b.api.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := handlers[i.ApplicationCommandData().Name]; ok {
			options := i.ApplicationCommandData().Options

			optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
			for _, opt := range options {
				optionMap[opt.Name] = opt
			}

			response := h(i, optionMap)
			_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: response,
				},
			})
		}
	})
}
