package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
)

// Bot parameters
var (
	GuildID        = flag.String("guild", "", "Test guild ID. If not passed - bot registers commands globally")
	BotToken       = flag.String("token", "", "Bot access token")
	RemoveCommands = flag.Bool("rmcmd", true, "Remove all commands after shutdowning or not")
	channelID      = flag.String("channel", "", "Channel ID")
)

var (
	dmPermission                   = false
	defaultMemberPermissions int64 = discordgo.PermissionManageServer

	commands = []*discordgo.ApplicationCommand{
		{
			Name:                     "permission-overview",
			Description:              "Command for demonstration of default command permissions",
			DefaultMemberPermissions: &defaultMemberPermissions,
			DMPermission:             &dmPermission,
		},
		{
			Name:        "random-image",
			Description: "get random image",
		},
		{
			Name:        "random-anime-image",
			Description: "get random anime image",
		},
		{
			Name:        "current-weather",
			Description: "get HCM city weather",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "city",
					Description: "input city name",
					Required:    true,
				},
			},
		},
		{
			Name:        "random-integer",
			Description: "get random integer",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "number",
					Description: "input max number you want to random",
					Required:    true,
				},
			},
		},
		{
			Name:        "lol-info",
			Description: "get lol info by name",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "name",
					Description: "input lol name",
					Required:    true,
				},
			},
		},
		{
			Name:        "lol-history",
			Description: "get lol history matches by name",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "name",
					Description: "input lol name",
					Required:    true,
				},
			},
		},
		{
			Name:        "pubg-info",
			Description: "get pubg info by name",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "name",
					Description: "input pubg steam name",
					Required:    true,
				},
			},
		},
		{
			Name:        "av-info",
			Description: "get av info by name",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "name",
					Description: "input av name",
					Required:    true,
				},
			},
		},
		{
			Name:        "subcommands",
			Description: "Subcommands and command groups example",
			Options: []*discordgo.ApplicationCommandOption{
				// When a command has subcommands/subcommand groups
				// It must not have top-level options, they aren't accesible in the UI
				// in this case (at least not yet), so if a command has
				// subcommands/subcommand any groups registering top-level options
				// will cause the registration of the command to fail

				{
					Name:        "subcommand-group",
					Description: "Subcommands group",
					Options: []*discordgo.ApplicationCommandOption{
						// Also, subcommand groups aren't capable of
						// containing options, by the name of them, you can see
						// they can only contain subcommands
						{
							Name:        "nested-subcommand",
							Description: "Nested subcommand",
							Type:        discordgo.ApplicationCommandOptionSubCommand,
						},
					},
					Type: discordgo.ApplicationCommandOptionSubCommandGroup,
				},
				// Also, you can create both subcommand groups and subcommands
				// in the command at the same time. But, there's some limits to
				// nesting, count of subcommands (top level and nested) and options.
				// Read the intro of slash-commands docs on Discord dev portal
				// to get more information
				{
					Name:        "subcommand",
					Description: "Top-level subcommand",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
				},
			},
		},
		{
			Name:        "responses",
			Description: "Interaction responses testing initiative",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "resp-type",
					Description: "Response type",
					Type:        discordgo.ApplicationCommandOptionInteger,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "Channel message with source",
							Value: 4,
						},
						{
							Name:  "Deferred response With Source",
							Value: 5,
						},
					},
					Required: true,
				},
			},
		},
		{
			Name:        "followups",
			Description: "Followup messages",
		},
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"random-image":       RandomImage(),
		"random-anime-image": RandomAnimeImage(),
		"current-weather":    TodayWeather(),
		"random-integer":     RandomInteger(),
		"lol-info":           LoLInfo(),
		"lol-history":        LoLMatches(),
		"pubg-info":          GetPUBGStats(),
		"av-info":            GetActress(),
		"permission-overview": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			perms, err := s.ApplicationCommandPermissions(s.State.User.ID, i.GuildID, i.ApplicationCommandData().ID)

			var restError *discordgo.RESTError
			if errors.As(err, &restError) && restError.Message != nil && restError.Message.Code == discordgo.ErrCodeUnknownApplicationCommandPermissions {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: ":x: No permission overwrites",
					},
				})
				return
			} else if err != nil {
				panic(err)
			}

			if err != nil {
				panic(err)
			}
			format := "- %s %s\n"

			channels := ""
			users := ""
			roles := ""

			for _, o := range perms.Permissions {
				emoji := "❌"
				if o.Permission {
					emoji = "☑"
				}

				switch o.Type {
				case discordgo.ApplicationCommandPermissionTypeUser:
					users += fmt.Sprintf(format, emoji, "<@!"+o.ID+">")
				case discordgo.ApplicationCommandPermissionTypeChannel:
					allChannels, _ := discordgo.GuildAllChannelsID(i.GuildID)

					if o.ID == allChannels {
						channels += fmt.Sprintf(format, emoji, "All channels")
					} else {
						channels += fmt.Sprintf(format, emoji, "<#"+o.ID+">")
					}
				case discordgo.ApplicationCommandPermissionTypeRole:
					if o.ID == i.GuildID {
						roles += fmt.Sprintf(format, emoji, "@everyone")
					} else {
						roles += fmt.Sprintf(format, emoji, "<@&"+o.ID+">")
					}
				}
			}

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{
						{
							Title:       "Permissions overview",
							Description: "Overview of permissions for this command",
							Fields: []*discordgo.MessageEmbedField{
								{
									Name:  "Users",
									Value: users,
								},
								{
									Name:  "Channels",
									Value: channels,
								},
								{
									Name:  "Roles",
									Value: roles,
								},
							},
						},
					},
					AllowedMentions: &discordgo.MessageAllowedMentions{},
				},
			})
		},
		"subcommands": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			options := i.ApplicationCommandData().Options
			content := ""

			// As you can see, names of subcommands (nested, top-level)
			// and subcommand groups are provided through the arguments.
			switch options[0].Name {
			case "subcommand":
				content = "The top-level subcommand is executed. Now try to execute the nested one."
			case "subcommand-group":
				options = options[0].Options
				switch options[0].Name {
				case "nested-subcommand":
					content = "Nice, now you know how to execute nested commands too"
				default:
					content = "Oops, something went wrong.\n" +
						"Hol' up, you aren't supposed to see this message."
				}
			}

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: content,
				},
			})
		},
		"responses": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			// Responses to a command are very important.
			// First of all, because you need to react to the interaction
			// by sending the response in 3 seconds after receiving, otherwise
			// interaction will be considered invalid and you can no longer
			// use the interaction token and ID for responding to the user's request

			content := ""
			// As you can see, the response type names used here are pretty self-explanatory,
			// but for those who want more information see the official documentation
			switch i.ApplicationCommandData().Options[0].IntValue() {
			case int64(discordgo.InteractionResponseChannelMessageWithSource):
				content =
					"You just responded to an interaction, sent a message and showed the original one. " +
						"Congratulations!"
				content +=
					"\nAlso... you can edit your response, wait 5 seconds and this message will be changed"
			default:
				err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseType(i.ApplicationCommandData().Options[0].IntValue()),
				})
				if err != nil {
					s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
						Content: "Something went wrong",
					})
				}
				return
			}

			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseType(i.ApplicationCommandData().Options[0].IntValue()),
				Data: &discordgo.InteractionResponseData{
					Content: content,
				},
			})
			if err != nil {
				s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
					Content: "Something went wrong",
				})
				return
			}
			time.AfterFunc(time.Second*5, func() {
				content := content + "\n\nWell, now you know how to create and edit responses. " +
					"But you still don't know how to delete them... so... wait 10 seconds and this " +
					"message will be deleted."
				_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
					Content: &content,
				})
				if err != nil {
					s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
						Content: "Something went wrong",
					})
					return
				}
				time.Sleep(time.Second * 10)
				s.InteractionResponseDelete(i.Interaction)
			})
		},
		"followups": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			// Followup messages are basically regular messages (you can create as many of them as you wish)
			// but work as they are created by webhooks and their functionality
			// is for handling additional messages after sending a response.

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					// Note: this isn't documented, but you can use that if you want to.
					// This flag just allows you to create messages visible only for the caller of the command
					// (user who triggered the command)
					Flags:   discordgo.MessageFlagsEphemeral,
					Content: "Surprise!",
				},
			})
			msg, err := s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
				Content: "Followup message has been created, after 5 seconds it will be edited",
			})
			if err != nil {
				s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
					Content: "Something went wrong",
				})
				return
			}
			time.Sleep(time.Second * 5)

			content := "Now the original message is gone and after 10 seconds this message will ~~self-destruct~~ be deleted."
			s.FollowupMessageEdit(i.Interaction, msg.ID, &discordgo.WebhookEdit{
				Content: &content,
			})

			time.Sleep(time.Second * 10)

			s.FollowupMessageDelete(i.Interaction, msg.ID)

			s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
				Content: "For those, who didn't skip anything and followed tutorial along fairly, " +
					"take a unicorn :unicorn: as reward!\n" +
					"Also, as bonus... look at the original interaction response :D",
			})
		},
	}
)

var s *discordgo.Session

func init() { flag.Parse() }

func init() {
	var err error
	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	if *BotToken == "" {
		*BotToken = os.Getenv("BOT_TOKEN")
	}
	if *channelID == "" {
		*channelID = os.Getenv("DISCORD_CRON_CHANNEL_ID")
	}

	s, err = discordgo.New("Bot " + *BotToken)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}
}

func init() {
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})

	//Cron
	c := cron.New()
	c.AddFunc("@every 12h0m0s", func() { HackerNews() })
	c.AddFunc("@every 12h0m0s", func() { CronGameNews() })
	c.AddFunc("@every 12h0m0s", func() { CronAvRecommend() })
	c.Start()
}

func main() {
	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})
	err := s.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}

	log.Println("Adding commands...")
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, *GuildID, v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}
	defer s.Close()

	//bot stop processing
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	log.Println("Press Ctrl+C to exit")
	<-stop

	if *RemoveCommands {
		log.Println("Removing commands...")
		// We need to fetch the commands, since deleting requires the command ID.
		// We are doing this from the returned commands on line 375, because using
		// this will delete all the commands, which might not be desirable, so we
		// are deleting only the commands that we added.
		registeredCommands, err := s.ApplicationCommands(s.State.User.ID, *GuildID)
		if err != nil {
			log.Fatalf("Could not fetch registered commands: %v", err)
		}

		for _, v := range registeredCommands {
			err := s.ApplicationCommandDelete(s.State.User.ID, *GuildID, v.ID)
			if err != nil {
				log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
			}
		}
	}
	log.Println("Gracefully shutting down.")
}

func HandleMessage(bot *discordgo.Session, message *discordgo.MessageCreate) {
	if message.Author.Bot || message.Author.ID == bot.State.User.ID {
		return
	}
}
