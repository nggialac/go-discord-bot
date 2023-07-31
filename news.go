package main

import (
	"encoding/json"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
)

var (
	HACKER_NEWS_TOP_STORIES_API  string = "https://hacker-news.firebaseio.com/v0/topstories.json"
	HACKER_NEWS_DETAIL_STORY_API string = "https://hacker-news.firebaseio.com/v0/item/"
	HACKER_NEWS_API_EXT          string = ".json"
)

var (
	NEWS_API_ROOT_URL           string = "https://newsapi.org/v2/"
	NEWS_API_TYPE_EVERYTHING    string = "everything"
	NEWS_API_TYPE_TOP_HEADLINES string = "top-headlines"
)

type DetailArticle struct {
	By          string `json:"by"`
	Descendants int    `json:"descendants"`
	ID          int    `json:"id"`
	Kids        []int  `json:"kids"`
	Score       int    `json:"score"`
	Time        int    `json:"time"`
	Title       string `json:"title"`
	Type        string `json:"type"`
	URL         string `json:"url"`
}

func GetTopStories() *DetailArticle {
	var stories []int
	s := ClientRequest(HACKER_NEWS_TOP_STORIES_API, nil, nil) // []int
	err := json.Unmarshal(s, &stories)
	if err != nil {
		log.Println("top stories error: ", err)
	}

	num := RandomInt(len(stories), 0)
	detailArticleURL := HACKER_NEWS_DETAIL_STORY_API + strconv.Itoa(stories[num]) + HACKER_NEWS_API_EXT
	article := ClientRequest(detailArticleURL, nil, nil)

	detailArticle := &DetailArticle{}
	err = json.Unmarshal(article, detailArticle)
	if err != nil {
		log.Println("top stories error: ", err)
	}

	return detailArticle
}

func HackerNews() func(*discordgo.Session, *discordgo.InteractionCreate) {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		article := GetTopStories()
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Title:   "Hacker News",
				Content: "Hacker News by " + article.By + "\n" + article.URL,
			},
		})
		if err != nil {
			log.Printf("get news error: " + err.Error())
			return
		}
	}
}

type NewsAPIResponse struct {
	Status       string `json:"status"`
	TotalResults int    `json:"totalResults"`
	Articles     []struct {
		Source struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"source"`
		Author      string    `json:"author"`
		Title       string    `json:"title"`
		Description string    `json:"description"`
		URL         string    `json:"url"`
		URLToImage  string    `json:"urlToImage"`
		PublishedAt time.Time `json:"publishedAt"`
		Content     string    `json:"content"`
	} `json:"articles"`
}

var (
	newsAPIParams = map[string]string{
		"page":     "1",
		"pageSize": "1",
		"q":        "football",
		"sortBy":   "publishedAt",
	}
)

func GetNewsTypeEverything(apiKey string) *NewsAPIResponse {
	url := NEWS_API_ROOT_URL + NEWS_API_TYPE_EVERYTHING
	newsAPIResponse := &NewsAPIResponse{}

	newAPIHeaders := map[string][]string{
		"x-api-key": {apiKey},
	}

	//Get yesterday's date
	yesterday := time.Now().Add(-24 * time.Hour).Format(time.DateOnly)
	if _, ok := newsAPIParams["from"]; !ok {
		newsAPIParams["from"] = yesterday
	}

	s := ClientRequest(url, newAPIHeaders, newsAPIParams) //
	err := json.Unmarshal(s, newsAPIResponse)
	if err != nil {
		log.Println("news api error: ", err)
	}

	return newsAPIResponse
}

func CronHackerNews() {
	article := GetTopStories()
	s.ChannelMessageSend(*channelID, "===News===\nAuthor: "+article.By+"\nTitle: "+article.Title+"\nLink: "+article.URL)
}

func CronGameNews() {
	apiKey := os.Getenv("NEWS_API_KEY")
	if apiKey == "" {
		log.Println("cannot get news api key")
		return
	}

	news := GetNewsTypeEverything(apiKey)
	if len(news.Articles) == 0 {
		log.Println("article not found")
		return
	}
	s.ChannelMessageSend(*channelID, "===News===\nAuthor: "+news.Articles[0].Author+"\nTitle: "+news.Articles[0].Title+"\nLink: "+news.Articles[0].URL)
}
