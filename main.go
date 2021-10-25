package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/tidwall/gjson"
)

// Variables used for command line parameters
var (
	Token string
)

const KuteGoAPIURL = "https://kutego-api-xxxxx-ew.a.run.app"
const WeatherAPIURL = "api.openweathermap.org/data/2.5/weather"
const WibuAPIURL = "https://api.waifu.pics"
const QuoteAPIURL = "https://api.kanye.rest"

func init() {
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

type Image struct {
	url string
}

func main() {
	// port := os.Getenv("PORT")

	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	// })

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + os.Getenv("BOT_TOKEN"))
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()

	// port := os.Getenv("PORT")
	// if port == "" {
	// 	port = "5000" // Default port if not specified
	// }

	// log.Fatal(http.ListenAndServe(":"+port, nil))
}

type Gopher struct {
	Name string `json: "name"`
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "!huybu" {
		jpg, err := os.Open("huybu.jpg")
		if err != nil {
			log.Fatal(err)
		}
		s.ChannelFileSendWithMessage(m.ChannelID, "Chuyên gia Phú Hòa!", "huybu.jpg", jpg)
	}

	if m.Content == "!hunghe" {
		jpg, err := os.Open("hunghe.jpg")
		if err != nil {
			log.Fatal(err)
		}
		s.ChannelFileSendWithMessage(m.ChannelID, "Chú hề Phú Lợi!", "hunghe.jpg", jpg)
	}

	if m.Content == "!nhan5" {
		jpg, err := os.Open("nhan5.jpg")
		if err != nil {
			log.Fatal(err)
		}
		s.ChannelFileSendWithMessage(m.ChannelID, "Mày muốn gì?", "nhan5.jpg", jpg)
	}

	if m.Content == "!minhpeo" {
		jpg, err := os.Open("minh.png")
		if err != nil {
			log.Fatal(err)
		}
		s.ChannelFileSendWithMessage(m.ChannelID, "Tôi kẹt giữa lòng sì gòn", "minh.png", jpg)
	}

	if m.Content == "!laclanhlung" {
		jpg, err := os.Open("lac.jpg")
		if err != nil {
			log.Fatal(err)
		}
		s.ChannelFileSendWithMessage(m.ChannelID, "Chúa tể đường giữa, ông vua chiến thuật, người cầm trịch trận đấu, Lạc số 2 không ai số 1!!!", "lac.jpg", jpg)
	}

	if m.Content == "!wibu" {
		response, err := http.Get(WibuAPIURL + "/sfw/waifu")
		if err != nil {
			fmt.Println(err)
		}

		defer response.Body.Close()
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(string(data))

		if response.StatusCode == 200 {
			value := gjson.Get(string(data), "url")
			println(value.String())
			_, err = s.ChannelMessageSend(m.ChannelID, value.String())
			if err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Println("Error: Can't get anime image! :-(")
		}
	}

	// if m.Content == "!wibulord" {
	// 	response, err := http.Get(WibuAPIURL + "/nsfw/waifu")
	// 	if err != nil {
	// 		fmt.Println(err)
	// 	}
	// 	defer response.Body.Close()
	// 	data, _ := ioutil.ReadAll(response.Body)
	// 	fmt.Println(string(data))

	// 	if response.StatusCode == 200 {
	// 		value := gjson.Get(string(data), "url")
	// 		println(value.String())
	// 		_, err = s.ChannelMessageSend(m.ChannelID, value.String())
	// 		if err != nil {
	// 			fmt.Println(err)
	// 		}
	// 	} else {
	// 		fmt.Println("Error: Can't get anime image! :-(")
	// 	}
	// }

	if m.Content == "!quote" {
		response, err := http.Get(QuoteAPIURL)
		if err != nil {
			fmt.Println(err)
		}
		defer response.Body.Close()
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(string(data))

		if response.StatusCode == 200 {
			value := gjson.Get(string(data), "quote")
			println(value.String())
			_, err = s.ChannelMessageSend(m.ChannelID, value.String())
			if err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Println("Error: Can't get quote! :-(")
		}
	}
}
