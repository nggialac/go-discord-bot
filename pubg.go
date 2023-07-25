package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/bwmarrin/discordgo"
)

// https://api.pubg.com/shards/steam/players?filter[playerNames]=Kio-sama
var (
	PUBGPlayerUrl string = "https://api.pubg.com/shards/steam/players/"
)

type PUBGPlayerReponse struct {
	Data []struct {
		Type       string `json:"type"`
		ID         string `json:"id"`
		Attributes struct {
			Name         string      `json:"name"`
			Stats        interface{} `json:"stats"`
			TitleID      string      `json:"titleId"`
			ShardID      string      `json:"shardId"`
			PatchVersion string      `json:"patchVersion"`
			BanType      string      `json:"banType"`
			ClanID       string      `json:"clanId"`
		} `json:"attributes"`
		Relationships struct {
			Assets struct {
				Data []interface{} `json:"data"`
			} `json:"assets"`
			Matches struct {
				Data []struct {
					Type string `json:"type"`
					ID   string `json:"id"`
				} `json:"data"`
			} `json:"matches"`
		} `json:"relationships"`
		Links struct {
			Self   string `json:"self"`
			Schema string `json:"schema"`
		} `json:"links"`
	} `json:"data"`
	Links struct {
		Self string `json:"self"`
	} `json:"links"`
	Meta struct {
	} `json:"meta"`
}

type PUBGPlayerStatsResponse struct {
	Data struct {
		Type       string `json:"type"`
		Attributes struct {
			GameModeStats struct {
				Duo struct {
					Assists             int     `json:"assists"`
					Boosts              int     `json:"boosts"`
					DBNOs               int     `json:"dBNOs"`
					DailyKills          int     `json:"dailyKills"`
					DailyWins           int     `json:"dailyWins"`
					DamageDealt         float64 `json:"damageDealt"`
					Days                int     `json:"days"`
					HeadshotKills       int     `json:"headshotKills"`
					Heals               int     `json:"heals"`
					KillPoints          int     `json:"killPoints"`
					Kills               int     `json:"kills"`
					LongestKill         float64 `json:"longestKill"`
					LongestTimeSurvived int     `json:"longestTimeSurvived"`
					Losses              int     `json:"losses"`
					MaxKillStreaks      int     `json:"maxKillStreaks"`
					MostSurvivalTime    int     `json:"mostSurvivalTime"`
					RankPoints          int     `json:"rankPoints"`
					RankPointsTitle     string  `json:"rankPointsTitle"`
					Revives             int     `json:"revives"`
					RideDistance        float64 `json:"rideDistance"`
					RoadKills           int     `json:"roadKills"`
					RoundMostKills      int     `json:"roundMostKills"`
					RoundsPlayed        int     `json:"roundsPlayed"`
					Suicides            int     `json:"suicides"`
					SwimDistance        float64 `json:"swimDistance"`
					TeamKills           int     `json:"teamKills"`
					TimeSurvived        float64 `json:"timeSurvived"`
					Top10S              int     `json:"top10s"`
					VehicleDestroys     int     `json:"vehicleDestroys"`
					WalkDistance        float64 `json:"walkDistance"`
					WeaponsAcquired     int     `json:"weaponsAcquired"`
					WeeklyKills         int     `json:"weeklyKills"`
					WeeklyWins          int     `json:"weeklyWins"`
					WinPoints           int     `json:"winPoints"`
					Wins                int     `json:"wins"`
				} `json:"duo"`
				DuoFpp struct {
					Assists             int     `json:"assists"`
					Boosts              int     `json:"boosts"`
					DBNOs               int     `json:"dBNOs"`
					DailyKills          int     `json:"dailyKills"`
					DailyWins           int     `json:"dailyWins"`
					DamageDealt         float64 `json:"damageDealt"`
					Days                int     `json:"days"`
					HeadshotKills       int     `json:"headshotKills"`
					Heals               int     `json:"heals"`
					KillPoints          int     `json:"killPoints"`
					Kills               int     `json:"kills"`
					LongestKill         float64 `json:"longestKill"`
					LongestTimeSurvived float64 `json:"longestTimeSurvived"`
					Losses              int     `json:"losses"`
					MaxKillStreaks      int     `json:"maxKillStreaks"`
					MostSurvivalTime    float64 `json:"mostSurvivalTime"`
					RankPoints          int     `json:"rankPoints"`
					RankPointsTitle     string  `json:"rankPointsTitle"`
					Revives             int     `json:"revives"`
					RideDistance        float64 `json:"rideDistance"`
					RoadKills           int     `json:"roadKills"`
					RoundMostKills      int     `json:"roundMostKills"`
					RoundsPlayed        int     `json:"roundsPlayed"`
					Suicides            int     `json:"suicides"`
					SwimDistance        int     `json:"swimDistance"`
					TeamKills           int     `json:"teamKills"`
					TimeSurvived        float64 `json:"timeSurvived"`
					Top10S              int     `json:"top10s"`
					VehicleDestroys     int     `json:"vehicleDestroys"`
					WalkDistance        float64 `json:"walkDistance"`
					WeaponsAcquired     int     `json:"weaponsAcquired"`
					WeeklyKills         int     `json:"weeklyKills"`
					WeeklyWins          int     `json:"weeklyWins"`
					WinPoints           int     `json:"winPoints"`
					Wins                int     `json:"wins"`
				} `json:"duo-fpp"`
				Solo struct {
					Assists             int     `json:"assists"`
					Boosts              int     `json:"boosts"`
					DBNOs               int     `json:"dBNOs"`
					DailyKills          int     `json:"dailyKills"`
					DailyWins           int     `json:"dailyWins"`
					DamageDealt         float64 `json:"damageDealt"`
					Days                int     `json:"days"`
					HeadshotKills       int     `json:"headshotKills"`
					Heals               int     `json:"heals"`
					KillPoints          int     `json:"killPoints"`
					Kills               int     `json:"kills"`
					LongestKill         float64 `json:"longestKill"`
					LongestTimeSurvived int     `json:"longestTimeSurvived"`
					Losses              int     `json:"losses"`
					MaxKillStreaks      int     `json:"maxKillStreaks"`
					MostSurvivalTime    int     `json:"mostSurvivalTime"`
					RankPoints          int     `json:"rankPoints"`
					RankPointsTitle     string  `json:"rankPointsTitle"`
					Revives             int     `json:"revives"`
					RideDistance        float64 `json:"rideDistance"`
					RoadKills           int     `json:"roadKills"`
					RoundMostKills      int     `json:"roundMostKills"`
					RoundsPlayed        int     `json:"roundsPlayed"`
					Suicides            int     `json:"suicides"`
					SwimDistance        float64 `json:"swimDistance"`
					TeamKills           int     `json:"teamKills"`
					TimeSurvived        int     `json:"timeSurvived"`
					Top10S              int     `json:"top10s"`
					VehicleDestroys     int     `json:"vehicleDestroys"`
					WalkDistance        float64 `json:"walkDistance"`
					WeaponsAcquired     int     `json:"weaponsAcquired"`
					WeeklyKills         int     `json:"weeklyKills"`
					WeeklyWins          int     `json:"weeklyWins"`
					WinPoints           int     `json:"winPoints"`
					Wins                int     `json:"wins"`
				} `json:"solo"`
				SoloFpp struct {
					Assists             int     `json:"assists"`
					Boosts              int     `json:"boosts"`
					DBNOs               int     `json:"dBNOs"`
					DailyKills          int     `json:"dailyKills"`
					DailyWins           int     `json:"dailyWins"`
					DamageDealt         float64 `json:"damageDealt"`
					Days                int     `json:"days"`
					HeadshotKills       int     `json:"headshotKills"`
					Heals               int     `json:"heals"`
					KillPoints          int     `json:"killPoints"`
					Kills               int     `json:"kills"`
					LongestKill         float64 `json:"longestKill"`
					LongestTimeSurvived int     `json:"longestTimeSurvived"`
					Losses              int     `json:"losses"`
					MaxKillStreaks      int     `json:"maxKillStreaks"`
					MostSurvivalTime    int     `json:"mostSurvivalTime"`
					RankPoints          int     `json:"rankPoints"`
					RankPointsTitle     string  `json:"rankPointsTitle"`
					Revives             int     `json:"revives"`
					RideDistance        float64 `json:"rideDistance"`
					RoadKills           int     `json:"roadKills"`
					RoundMostKills      int     `json:"roundMostKills"`
					RoundsPlayed        int     `json:"roundsPlayed"`
					Suicides            int     `json:"suicides"`
					SwimDistance        int     `json:"swimDistance"`
					TeamKills           int     `json:"teamKills"`
					TimeSurvived        float64 `json:"timeSurvived"`
					Top10S              int     `json:"top10s"`
					VehicleDestroys     int     `json:"vehicleDestroys"`
					WalkDistance        float64 `json:"walkDistance"`
					WeaponsAcquired     int     `json:"weaponsAcquired"`
					WeeklyKills         int     `json:"weeklyKills"`
					WeeklyWins          int     `json:"weeklyWins"`
					WinPoints           int     `json:"winPoints"`
					Wins                int     `json:"wins"`
				} `json:"solo-fpp"`
				Squad struct {
					Assists             int     `json:"assists"`
					Boosts              int     `json:"boosts"`
					DBNOs               int     `json:"dBNOs"`
					DailyKills          int     `json:"dailyKills"`
					DailyWins           int     `json:"dailyWins"`
					DamageDealt         float64 `json:"damageDealt"`
					Days                int     `json:"days"`
					HeadshotKills       int     `json:"headshotKills"`
					Heals               int     `json:"heals"`
					KillPoints          int     `json:"killPoints"`
					Kills               int     `json:"kills"`
					LongestKill         float64 `json:"longestKill"`
					LongestTimeSurvived int     `json:"longestTimeSurvived"`
					Losses              int     `json:"losses"`
					MaxKillStreaks      int     `json:"maxKillStreaks"`
					MostSurvivalTime    int     `json:"mostSurvivalTime"`
					RankPoints          int     `json:"rankPoints"`
					RankPointsTitle     string  `json:"rankPointsTitle"`
					Revives             int     `json:"revives"`
					RideDistance        float64 `json:"rideDistance"`
					RoadKills           int     `json:"roadKills"`
					RoundMostKills      int     `json:"roundMostKills"`
					RoundsPlayed        int     `json:"roundsPlayed"`
					Suicides            int     `json:"suicides"`
					SwimDistance        float64 `json:"swimDistance"`
					TeamKills           int     `json:"teamKills"`
					TimeSurvived        float64 `json:"timeSurvived"`
					Top10S              int     `json:"top10s"`
					VehicleDestroys     int     `json:"vehicleDestroys"`
					WalkDistance        float64 `json:"walkDistance"`
					WeaponsAcquired     int     `json:"weaponsAcquired"`
					WeeklyKills         int     `json:"weeklyKills"`
					WeeklyWins          int     `json:"weeklyWins"`
					WinPoints           int     `json:"winPoints"`
					Wins                int     `json:"wins"`
				} `json:"squad"`
				SquadFpp struct {
					Assists             int     `json:"assists"`
					Boosts              int     `json:"boosts"`
					DBNOs               int     `json:"dBNOs"`
					DailyKills          int     `json:"dailyKills"`
					DailyWins           int     `json:"dailyWins"`
					DamageDealt         float64 `json:"damageDealt"`
					Days                int     `json:"days"`
					HeadshotKills       int     `json:"headshotKills"`
					Heals               int     `json:"heals"`
					KillPoints          int     `json:"killPoints"`
					Kills               int     `json:"kills"`
					LongestKill         float64 `json:"longestKill"`
					LongestTimeSurvived int     `json:"longestTimeSurvived"`
					Losses              int     `json:"losses"`
					MaxKillStreaks      int     `json:"maxKillStreaks"`
					MostSurvivalTime    int     `json:"mostSurvivalTime"`
					RankPoints          int     `json:"rankPoints"`
					RankPointsTitle     string  `json:"rankPointsTitle"`
					Revives             int     `json:"revives"`
					RideDistance        float64 `json:"rideDistance"`
					RoadKills           int     `json:"roadKills"`
					RoundMostKills      int     `json:"roundMostKills"`
					RoundsPlayed        int     `json:"roundsPlayed"`
					Suicides            int     `json:"suicides"`
					SwimDistance        float64 `json:"swimDistance"`
					TeamKills           int     `json:"teamKills"`
					TimeSurvived        int     `json:"timeSurvived"`
					Top10S              int     `json:"top10s"`
					VehicleDestroys     int     `json:"vehicleDestroys"`
					WalkDistance        float64 `json:"walkDistance"`
					WeaponsAcquired     int     `json:"weaponsAcquired"`
					WeeklyKills         int     `json:"weeklyKills"`
					WeeklyWins          int     `json:"weeklyWins"`
					WinPoints           int     `json:"winPoints"`
					Wins                int     `json:"wins"`
				} `json:"squad-fpp"`
			} `json:"gameModeStats"`
			BestRankPoint int `json:"bestRankPoint"`
		} `json:"attributes"`
		Relationships struct {
			MatchesSquad struct {
				Data []struct {
					Type string `json:"type"`
					ID   string `json:"id"`
				} `json:"data"`
			} `json:"matchesSquad"`
			MatchesSquadFPP struct {
				Data []struct {
					Type string `json:"type"`
					ID   string `json:"id"`
				} `json:"data"`
			} `json:"matchesSquadFPP"`
			Season struct {
				Data struct {
					Type string `json:"type"`
					ID   string `json:"id"`
				} `json:"data"`
			} `json:"season"`
			Player struct {
				Data struct {
					Type string `json:"type"`
					ID   string `json:"id"`
				} `json:"data"`
			} `json:"player"`
			MatchesSolo struct {
				Data []interface{} `json:"data"`
			} `json:"matchesSolo"`
			MatchesSoloFPP struct {
				Data []struct {
					Type string `json:"type"`
					ID   string `json:"id"`
				} `json:"data"`
			} `json:"matchesSoloFPP"`
			MatchesDuo struct {
				Data []interface{} `json:"data"`
			} `json:"matchesDuo"`
			MatchesDuoFPP struct {
				Data []interface{} `json:"data"`
			} `json:"matchesDuoFPP"`
		} `json:"relationships"`
	} `json:"data"`
	Links struct {
		Self string `json:"self"`
	} `json:"links"`
	Meta struct {
	} `json:"meta"`
}

func FindPUBGPlayer(name string) (*PUBGPlayerStatsResponse, error) {
	// region := "pc-sea"
	// gameMode := "squad"
	// seasonID := "division.bro.official.pc-2018-24"

	res := &PUBGPlayerReponse{}
	apiID := os.Getenv("PUBG_API_KEY")

	var headers = map[string][]string{
		"Accept":        {"application/vnd.api+json"},
		"Authorization": {apiID},
	}

	var params = map[string]string{
		"filter[playerNames]": name,
	}

	playerData := ClientRequest(PUBGPlayerUrl, headers, params)
	err := json.Unmarshal(playerData, res)
	if err != nil {
		log.Println("find player error: ", err)
		return nil, err
	}

	var playerID string
	if len(res.Data) == 0 {
		return nil, fmt.Errorf("Cannot find player")
	}
	playerID = res.Data[0].ID

	lifetimeStatsURL := playerID + "/seasons/lifetime"
	statsRes := &PUBGPlayerStatsResponse{}
	statsData := ClientRequest(PUBGPlayerUrl+lifetimeStatsURL, headers, nil)
	err = json.Unmarshal(statsData, statsRes)
	if err != nil {
		log.Println("find stats error: ", err)
		return nil, err
	}

	return statsRes, err
}

func GetPUBGStats() func(*discordgo.Session, *discordgo.InteractionCreate) {
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

		playerStats, err := FindPUBGPlayer(name)
		if err != nil {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Cannot find this account, err: " + name,
				},
			})
			return
		}
		squad := playerStats.Data.Attributes.GameModeStats.Squad

		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Account name: " + name,
				Embeds: []*discordgo.MessageEmbed{
					{
						Title: "SQUAD Mode",
						Image: &discordgo.MessageEmbedImage{URL: "https://images2.thanhnien.vn/Uploaded/duongntt/2020_07_16/pubg-mobile-huong-dan-su-dung-cac-phu-kien-hieu-qua-nhat1542009065_VJZV.jpg?width=500"},
						Fields: []*discordgo.MessageEmbedField{
							{
								Name:  "Kills",
								Value: strconv.Itoa(squad.Kills),
							},
							{
								Name:  "Assists",
								Value: strconv.Itoa(squad.Assists),
							},
							{
								Name:  "Wins",
								Value: strconv.Itoa(squad.Wins),
							},
							{
								Name:  "Damage dealt",
								Value: strconv.Itoa(int(squad.DamageDealt)),
							},
							{
								Name:  "Heals",
								Value: strconv.Itoa(squad.Heals),
							},
							{
								Name:  "Headshot kills",
								Value: strconv.Itoa(squad.HeadshotKills),
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
