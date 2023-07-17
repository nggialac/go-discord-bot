package main

import (
	"log"
	"strconv"

	"fmt"

	"github.com/bwmarrin/discordgo"
)

var (
	images = []*discordgo.MessageEmbed{
		{Image: &discordgo.MessageEmbedImage{URL: "https://scontent.fsgn4-1.fna.fbcdn.net/v/t1.15752-9/361114965_1377841329614903_3059549030702874547_n.jpg?_nc_cat=103&cb=99be929b-59f725be&ccb=1-7&_nc_sid=ae9488&_nc_ohc=NuxE2MWSfdYAX9SIx4S&_nc_ht=scontent.fsgn4-1.fna&oh=03_AdRA2s-70q-rsoCsN8LyxzwjxSxqIackDZ9z0UyrMl87Hg&oe=64D9D33F"}},
		{Image: &discordgo.MessageEmbedImage{URL: "https://scontent.fsgn4-1.fna.fbcdn.net/v/t1.15752-9/355395717_288793423702023_9066634740602676872_n.jpg?_nc_cat=103&cb=99be929b-59f725be&ccb=1-7&_nc_sid=ae9488&_nc_ohc=RZlLGkUxUgIAX8yLQR5&_nc_ht=scontent.fsgn4-1.fna&oh=03_AdRL2W5KVMux_F5hd_9Eefy2aXTvoF_XZzcqpng3lh9bQQ&oe=64D9B435"}},
		{Image: &discordgo.MessageEmbedImage{URL: "https://scontent.fsgn4-1.fna.fbcdn.net/v/t1.15752-9/353503087_929245591700525_5176460622788488829_n.jpg?_nc_cat=101&cb=99be929b-59f725be&ccb=1-7&_nc_sid=ae9488&_nc_ohc=uzvE1tYonK4AX8MsW5j&_nc_ht=scontent.fsgn4-1.fna&oh=03_AdRdN26f3wZBsECURul7X2qjLuSvlsUnAC7FD_-kNwNKTw&oe=64D9D408"}},
	}
)

func RandomImage() func(*discordgo.Session, *discordgo.InteractionCreate) {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		imageIndex := RandomInt(len(images), 0) //random image index
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "This is random image",
				Embeds: []*discordgo.MessageEmbed{
					images[imageIndex],
				},
			},
		})

		if err != nil {
			log.Printf("random image error: " + err.Error())
		}
	}
}

func RandomAnimeImage() func(*discordgo.Session, *discordgo.InteractionCreate) {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "This is random image",
				Embeds: []*discordgo.MessageEmbed{
					{Image: &discordgo.MessageEmbedImage{URL: GetAnimeImage()}},
				},
			},
		})

		if err != nil {
			log.Printf("random anime image error: " + err.Error())
		}
	}
}

func TodayWeather() func(*discordgo.Session, *discordgo.InteractionCreate) {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		options := i.ApplicationCommandData().Options
		optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
		for _, opt := range options {
			optionMap[opt.Name] = opt
		}

		var city string
		if option, ok := optionMap["city"]; ok {
			city = option.StringValue()
		}

		c := GetCurrentWeather(city)
		icon := "https://" + c.Current.Condition.Icon[2:]
		log.Print("Icon: ", icon)

		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Current weather: ",
				Embeds: []*discordgo.MessageEmbed{
					// {Image: &discordgo.MessageEmbedImage{URL: c.Current.Condition.Icon}},
					{
						Title:       c.Location.Name + ", " + c.Location.Country,
						Description: "CONDITION: " + c.Current.Condition.Text,
						Thumbnail:   &discordgo.MessageEmbedThumbnail{URL: icon},
						Fields: []*discordgo.MessageEmbedField{
							{
								Name:  "Last updated",
								Value: c.Current.LastUpdated,
							},
							{
								Name:  "Humidity",
								Value: strconv.Itoa(c.Current.Humidity),
							},
							{
								Name:  "Cloud",
								Value: strconv.Itoa(c.Current.Cloud),
							},
							{
								Name:  "Temp C",
								Value: fmt.Sprintf("%f", c.Current.TempC),
							},
							{
								Name:  "UV",
								Value: fmt.Sprintf("%f", c.Current.Uv),
							},
						},
					},
				},
			},
		})
		if err != nil {
			log.Printf("weather error: %s" + err.Error())
		}
	}
}

func RandomInteger() func(*discordgo.Session, *discordgo.InteractionCreate) {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		options := i.ApplicationCommandData().Options
		optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
		for _, opt := range options {
			optionMap[opt.Name] = opt
		}

		var number int64
		if option, ok := optionMap["number"]; ok {
			number = option.IntValue()
		}

		res := RandomInt(int(number)+1, 1)
		pokeID := RandomInt(500, 1)
		pokemonInfo := GetPokemonImage(pokeID)

		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "This is result: " + strconv.Itoa(res),
				Embeds: []*discordgo.MessageEmbed{
					{
						Title: pokemonInfo.Name,
						Image: &discordgo.MessageEmbedImage{URL: pokemonInfo.ImageUrl}},
				},
			},
		})

		if err != nil {
			log.Printf("random error: " + err.Error())
		}
	}
}
