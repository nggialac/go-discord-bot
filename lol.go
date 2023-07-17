package main

import (
	"fmt"
	"log"
	"os"

	"github.com/KnutZuidema/golio"
	"github.com/KnutZuidema/golio/api"
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

type LoLRankImage string

const (
	IRON_IMAGE_URL        LoLRankImage = "https://i.pinimg.com/564x/f5/a5/9d/f5a59d542851a7dcb3d0eae1851af735.jpg"
	BRONZE_IMAGE_URL      LoLRankImage = "https://i.pinimg.com/564x/bf/4d/40/bf4d401aa059647cc24cc8408f203e44.jpg"
	SILVER_IMAGE_URL      LoLRankImage = "https://i.pinimg.com/564x/75/61/5a/75615a37309f44c6f07353277429a4f2.jpg"
	GOLD_IMAGE_URL        LoLRankImage = "https://i.pinimg.com/564x/d7/58/1b/d7581b2a1033309523d20c9d1a1f4589.jpg"
	PLATINUM_IMAGE_URL    LoLRankImage = "https://i.pinimg.com/564x/d7/47/1e/d7471e2ef48175986e9b75b566f61408.jpg"
	DIAMOND_IMAGE_URL     LoLRankImage = "https://i.pinimg.com/564x/6a/10/c7/6a10c7e84c9f4e4aa9412582d28f3fd2.jpg"
	MASTER_IMAGE_URL      LoLRankImage = "https://i.pinimg.com/564x/69/61/ab/6961ab1af799f02df28fa74278d78120.jpg"
	GRANDMASTER_IMAGE_URL LoLRankImage = "https://i.pinimg.com/564x/ae/dd/2c/aedd2c30290af7f2a9b343027b31b0d2.jpg"
	CHALLENGER_IMAGE_URL  LoLRankImage = "https://i.pinimg.com/564x/b5/94/2f/b5942fb47954ab7756edceea90c7b052.jpg"
)

var (
	ASSETS_PATH                              = "./assets/"
	ASSETS_IMAGE_EXT                         = ".png"
	LoLRankingName   map[string]LoLRankImage = map[string]LoLRankImage{
		"BRONZE":       IRON_IMAGE_URL,
		"SILVER":       SILVER_IMAGE_URL,
		"GOLD":         GOLD_IMAGE_URL,
		"PLATINUM":     PLATINUM_IMAGE_URL,
		"DIAMOND":      DIAMOND_IMAGE_URL,
		"MASTER":       MASTER_IMAGE_URL,
		"GRAND MASTER": GRANDMASTER_IMAGE_URL,
		"CHALLENGER":   CHALLENGER_IMAGE_URL,
	}
)

func LOLClient() *golio.Client {
	return golio.NewClient(
		os.Getenv("RIOT_TOKEN"),
		golio.WithRegion(api.RegionVietNam),
		golio.WithLogger(logrus.New()))
}

func LOLInfo() func(*discordgo.Session, *discordgo.InteractionCreate) {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		options := i.ApplicationCommandData().Options
		optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
		for _, opt := range options {
			optionMap[opt.Name] = opt
		}
		var name string
		if option, ok := optionMap["name"]; ok {
			name = option.StringValue()
		}
		client := LOLClient()

		//Get account by name (SEA region)
		summoner, err := client.Riot.LoL.Summoner.GetByName(name)
		if err != nil {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Cannot find this account with name: " + name,
				},
			})
		}
		id := summoner.ID

		leagueItems, _ := client.Riot.LoL.League.ListBySummoner(id)
		normalRankInfo := leagueItems[0]
		tftRankInfo := leagueItems[1]

		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Account name: " + name,
				Embeds: []*discordgo.MessageEmbed{
					{
						Title: "Normal 5vs5",
						Image: &discordgo.MessageEmbedImage{URL: string(LoLRankingName[normalRankInfo.Tier])},
						Fields: []*discordgo.MessageEmbedField{{
							Name:  "Tier",
							Value: normalRankInfo.Tier,
						},
							{
								Name:  "Rank",
								Value: normalRankInfo.Rank,
							},
							{
								Name:  "League Points",
								Value: fmt.Sprintf("%d", normalRankInfo.LeaguePoints),
							},
							{
								Name:  "Ratio",
								Value: fmt.Sprintf("%d/%d", normalRankInfo.Wins, normalRankInfo.Losses),
							},
						},
					},
					{
						Title: "Team Tactics Fight",
						Image: &discordgo.MessageEmbedImage{URL: string(LoLRankingName[tftRankInfo.Tier])},
						Fields: []*discordgo.MessageEmbedField{{
							Name:  "Tier",
							Value: tftRankInfo.Tier,
						},
							{
								Name:  "Rank",
								Value: tftRankInfo.Rank,
							},
							{
								Name:  "League Points",
								Value: fmt.Sprintf("%d", tftRankInfo.LeaguePoints),
							},
							{
								Name:  "Ratio",
								Value: fmt.Sprintf("%d/%d", tftRankInfo.Wins, tftRankInfo.Losses),
							},
						},
					},
				},
			},
		})
		if err != nil {
			log.Printf("get rank error: " + err.Error())
		}
	}
}
