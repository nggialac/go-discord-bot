package main

import (
	"encoding/json"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/gojp/kana"
)

var (
	AV_API_URL         string = "https://api.dmm.com/affiliate/v3/ActressSearch"
	AV_PRODUCT_API_URL string = "https://api.dmm.com/affiliate/v3/ItemList"
	AV_SERIAL_URL      string = "https://api.dmm.com/affiliate/v3/SeriesSearch"
)

type AvAPIResponse struct {
	Request struct {
		Parameters struct {
			APIID       string `json:"api_id"`
			AffiliateID string `json:"affiliate_id"`
			Output      string `json:"output"`
			Offset      string `json:"offset"`
			Hits        string `json:"hits"`
			Keyword     string `json:"keyword"`
		} `json:"parameters"`
	} `json:"request"`
	Result struct {
		Status        string `json:"status"`
		ResultCount   int    `json:"result_count"`
		TotalCount    string `json:"total_count"`
		FirstPosition string `json:"first_position"`
		Actress       []struct {
			ID          string `json:"id"`
			Name        string `json:"name,omitempty"`
			Ruby        string `json:"ruby,omitempty"`
			Bust        string `json:"bust,omitempty"`
			Cup         string `json:"cup,omitempty"`
			Waist       string `json:"waist,omitempty"`
			Hip         string `json:"hip,omitempty"`
			Height      string `json:"height,omitempty"`
			Birthday    string `json:"birthday,omitempty"`
			BloodType   string `json:"blood_type,omitempty"`
			Hobby       string `json:"hobby,omitempty"`
			Prefectures string `json:"prefectures,omitempty"`
			ImageURL    struct {
				Small string `json:"small,omitempty"`
				Large string `json:"large,omitempty"`
			} `json:"imageURL"`
			ListURL struct {
				Digital string `json:"digital,omitempty"`
				Monthly string `json:"monthly,omitempty"`
				Mono    string `json:"mono,omitempty"`
				Rental  string `json:"rental,omitempty"`
			} `json:"listURL"`
		} `json:"actress"`
	} `json:"result"`
}

func FindActress(name string) *AvAPIResponse {
	res := &AvAPIResponse{}
	hiraganaName := kana.RomajiToHiragana(strings.ToLower(name))
	apiID := os.Getenv("AV_API_ID")

	var params = map[string]string{
		"output":       "json",
		"offset":       "1",
		"hits":         "10",
		"api_id":       apiID,
		"affiliate_id": "10278-996",
		"keyword":      hiraganaName,
	}

	data := ClientRequest(AV_API_URL, nil, params) // []int
	err := json.Unmarshal(data, res)
	if err != nil {
		log.Println("find actress error: ", err)
	}

	return res
}

func GetActress() func(s *discordgo.Session, i *discordgo.InteractionCreate) {
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

		res := FindActress(name)
		if res.Result.ResultCount == 0 {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Cannot find av: " + name,
				},
			})
			return
		}
		actress := res.Result.Actress[0]

		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Actress: " + actress.Name,
				Embeds: []*discordgo.MessageEmbed{
					{
						Title: actress.Ruby,
						Image: &discordgo.MessageEmbedImage{URL: actress.ImageURL.Large},
						Fields: []*discordgo.MessageEmbedField{
							{
								Name:  "Bust",
								Value: actress.Bust,
							},
							{
								Name:  "Cup",
								Value: actress.Cup,
							},
							{
								Name:  "Waist",
								Value: actress.Waist,
							},
							{
								Name:  "Hip",
								Value: actress.Hip,
							},
							{
								Name:  "Height",
								Value: actress.Height,
							},
							{
								Name:  "Birthday",
								Value: actress.Birthday,
							},
						},
					},
				},
			},
		})

		if err != nil {
			log.Printf("random error: " + err.Error())
		}
	}
}

func GetRecommendActress() (*AvAPIResponse, error) {
	apiID := os.Getenv("AV_API_ID")
	offset := RandomInt(1000, 1)

	var params = map[string]string{
		"service":      "digital",
		"output":       "json",
		"offset":       strconv.Itoa(offset),
		"hits":         "20",
		"api_id":       apiID,
		"affiliate_id": "10278-996",
		"article":      "actress",
		"gte_birthday": time.Now().AddDate(-35, 0, 0).Format("2006-01-02"),
	}

	data := ClientRequest(AV_API_URL, nil, params)

	res := &AvAPIResponse{}
	err := json.Unmarshal(data, res)
	if err != nil {
		log.Println("av error: ", err)
		return nil, err
	}

	if res.Result.ResultCount == 0 {
		log.Println("retry: ", err)
		return nil, err
	}

	return res, nil
}

func GetRecommendCode() (*AvAPIResponse, error) {
	apiID := os.Getenv("AV_API_ID")
	offset := RandomInt(1000, 1)

	var params = map[string]string{
		"service":      "digital",
		"output":       "json",
		"offset":       strconv.Itoa(offset),
		"hits":         "20",
		"api_id":       apiID,
		"affiliate_id": "10278-996",
		"article":      "actress",
		"gte_birthday": time.Now().AddDate(-35, 0, 0).Format("2006-01-02"),
	}

	data := ClientRequest(AV_API_URL, nil, params)

	res := &AvAPIResponse{}
	err := json.Unmarshal(data, res)
	if err != nil {
		log.Println("av error: ", err)
		return nil, err
	}

	if res.Result.ResultCount == 0 {
		log.Println("retry: ", err)
		return nil, err
	}

	return res, nil
}

func CronAvRecommend() {
	privateChannelID := os.Getenv("DISCORD_CRON_PRIVATE_CHANNEL_ID")
	if privateChannelID == "" {
		log.Fatal("private channel ID is wrong")
	}

	res, err := GetRecommendActress()
	if err != nil {
		log.Println("get actress was error: ", err)
		s.ChannelMessageSend(privateChannelID, "===AV===\nActress got an error")
		return
	}

	index := RandomInt(res.Result.ResultCount-1, 0)
	a := res.Result.Actress[index]

	_, err = s.ChannelMessageSendEmbed(privateChannelID, &discordgo.MessageEmbed{
		Title:       "ID " + a.ID,
		Description: "Name: " + a.Name + ", Ruby: " + a.Ruby,
		Thumbnail:   &discordgo.MessageEmbedThumbnail{URL: a.ImageURL.Large},
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "Bust",
				Value: a.Bust,
			},
			{
				Name:   "Cup",
				Value:  a.Cup,
				Inline: true,
			},
			{
				Name:   "Waist",
				Value:  a.Waist,
				Inline: true,
			},
			{
				Name:   "Hip",
				Value:  a.Hip,
				Inline: true,
			},
			{
				Name:  "Height",
				Value: a.Height,
			},
			{
				Name:   "Birthday",
				Value:  a.Birthday,
				Inline: true,
			},
			{
				Name:   "Blood Type",
				Value:  a.BloodType,
				Inline: true,
			},
			{
				Name:   "Height",
				Value:  a.Height,
				Inline: true,
			},
			{
				Name:  "Digital",
				Value: a.ListURL.Digital,
			},
			{
				Name:  "Monthly",
				Value: a.ListURL.Monthly,
			},
		},
	})
	if err != nil {
		log.Println("get actress was error: ", err)
		s.ChannelMessageSend(privateChannelID, "===AV===\nActress got an error")
		return
	}
}
