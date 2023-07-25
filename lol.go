package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"sync"

	"github.com/KnutZuidema/golio"
	"github.com/KnutZuidema/golio/api"
	"github.com/KnutZuidema/golio/riot/lol"
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

	CHAMPION_IMAGE_ROOT_URL string = "https://ddragon.leagueoflegends.com/cdn/12.4.1/img/champion"
	CHAMPION_IMAGE_EXT      string = "png"
)

var (
	LoLRankingName map[string]LoLRankImage = map[string]LoLRankImage{
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

func LoLClient() *golio.Client {
	return golio.NewClient(
		os.Getenv("RIOT_TOKEN"),
		golio.WithRegion(api.RegionVietnam),
		golio.WithLogger(logrus.New()))
}

func LoLInfo() func(*discordgo.Session, *discordgo.InteractionCreate) {
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
		client := LoLClient()

		//Get account by name (SEA region)
		summoner, err := client.Riot.LoL.Summoner.GetByName(name)
		if err != nil {
			log.Println("get name was error: ", err)
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Cannot find this account with name: " + name,
				},
			})
			return
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
						Title: "Rank Solo 5vs5",
						Image: &discordgo.MessageEmbedImage{URL: string(LoLRankingName[normalRankInfo.Tier])},
						Fields: []*discordgo.MessageEmbedField{
							{
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
						Fields: []*discordgo.MessageEmbedField{
							{
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
			return
		}
	}
}

func LoLMatches() func(*discordgo.Session, *discordgo.InteractionCreate) {
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
		client := LoLClient()

		//Get account by name (SEA region)
		summoner, err := client.Riot.LoL.Summoner.GetByName(name)
		if err != nil {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Cannot find this account with name: " + name,
				},
			})
			return
		}
		puuid := summoner.PUUID

		//Get matches - 5
		LoLMatches, err := client.Riot.LoL.Match.List(puuid, 0, 3, &lol.MatchListOptions{})
		if err != nil {
			fmt.Println("Cannot get matches id", err)
		}

		//Details
		matches := AsyncGetMatches(client, LoLMatches)
		sort.Slice(matches, func(i, j int) bool {
			return matches[i].Info.GameStartTimestamp > matches[j].Info.GameStartTimestamp
		})

		//search details by puuid
		var details []*lol.Participant
		for i := 0; i < len(matches); i++ {
			for _, v := range matches[i].Info.Participants {
				if v.PUUID == puuid {
					details = append(details, v)
				}
			}
		}

		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Account name: " + name,
				Embeds: []*discordgo.MessageEmbed{
					{
						Title: "Champ: " + details[0].ChampionName,
						Image: &discordgo.MessageEmbedImage{URL: GetChampionImageUrl(details[0].ChampionName)},
						Fields: []*discordgo.MessageEmbedField{
							{
								Name:  "Time Start",
								Value: ConvertUnixTime(matches[0].Info.GameStartTimestamp),
							},
							{
								Name:  "Mode",
								Value: matches[0].Info.GameMode,
							},
							{
								Name:  "Duration",
								Value: SecondsToMinutes(matches[0].Info.GameDuration),
							},
							{
								Name:  "K/D/A",
								Value: fmt.Sprintf("%d/%d/%d", details[0].Kills, details[0].Deaths, details[0].Assists),
							},
							{
								Name:  "DamageTaken/DamageDealt",
								Value: fmt.Sprintf("%d/%d", details[0].TotalDamageTaken, details[0].TotalDamageDealt),
							},
							{
								Name:  "Is Win?",
								Value: strconv.FormatBool(details[0].Win),
							},
						},
					},
					{
						Title: "Champ: " + details[1].ChampionName,
						Image: &discordgo.MessageEmbedImage{URL: GetChampionImageUrl(details[1].ChampionName)},
						Fields: []*discordgo.MessageEmbedField{
							{
								Name:  "Time Start",
								Value: ConvertUnixTime(matches[1].Info.GameStartTimestamp),
							},
							{
								Name:  "Mode",
								Value: matches[1].Info.GameMode,
							},
							{
								Name:  "Duration",
								Value: SecondsToMinutes(matches[1].Info.GameDuration),
							},
							{
								Name:  "K/D/A",
								Value: fmt.Sprintf("%d/%d/%d", details[1].Kills, details[1].Deaths, details[1].Assists),
							},
							{
								Name:  "DamageTaken/DamageDealt",
								Value: fmt.Sprintf("%d/%d", details[1].TotalDamageTaken, details[1].TotalDamageDealt),
							},
							{
								Name:  "Is Win?",
								Value: strconv.FormatBool(details[1].Win),
							},
						},
					},
					{
						Title: "Champ: " + details[2].ChampionName,
						Image: &discordgo.MessageEmbedImage{URL: GetChampionImageUrl(details[2].ChampionName)},
						Fields: []*discordgo.MessageEmbedField{
							{
								Name:  "Time Start",
								Value: ConvertUnixTime(matches[2].Info.GameStartTimestamp),
							},
							{
								Name:  "Mode",
								Value: matches[2].Info.GameMode,
							},
							{
								Name:  "Duration",
								Value: SecondsToMinutes(matches[2].Info.GameDuration),
							},
							{
								Name:  "K/D/A",
								Value: fmt.Sprintf("%d/%d/%d", details[2].Kills, details[2].Deaths, details[2].Assists),
							},
							{
								Name:  "DamageTaken/DamageDealt",
								Value: fmt.Sprintf("%d/%d", details[2].TotalDamageTaken, details[2].TotalDamageDealt),
							},
							{
								Name:  "Is Win?",
								Value: strconv.FormatBool(details[2].Win),
							},
						},
					},
				},
			},
		})
		if err != nil {
			log.Printf("get account history error: " + err.Error())
			return
		}
	}
}

func AsyncGetMatches(c *golio.Client, LoLMatches []string) []*lol.Match {
	var matches []*lol.Match
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for i := 0; i < len(LoLMatches); i++ {
			MatchStats, _ := c.Riot.LoL.Match.Get(LoLMatches[i])
			matches = append(matches, MatchStats)
		}
		wg.Done()
	}()
	wg.Wait()

	return matches
}

func SecondsToMinutes(inSeconds int) string {
	minutes := inSeconds / 60
	seconds := inSeconds % 60
	str := fmt.Sprintf("%d:%d", minutes, seconds)
	return str
}

func GetChampionImageUrl(champ string) string {
	return fmt.Sprintf("%s/%s.%s", CHAMPION_IMAGE_ROOT_URL, champ, CHAMPION_IMAGE_EXT)
}
