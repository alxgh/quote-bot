package main

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	
	"os"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"fmt"
	"strings"
)

const (
	QUOTE_URL = "http://quotesondesign.com/wp-json/posts?filter[orderby]=rand&filter[posts_per_page]=1"
)


type Quote struct {
	ID int 
	Title string
	Content string
	Link string 
}


func message(input string) string {
	if input == "help" || input == "/help" {
		return "There is some command: \n\n /quote -> Show a quote."
	}

	if input == "start" || input == "/start" {
		return "Welcome to our bot!"
	}

	if input == "/quote" || input == "quote" {
		resp, err := http.Get(QUOTE_URL)
		if err != nil {
			panic(err)
		}

		defer resp.Body.Close()
		
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		
		var quotedata  []Quote
		json.Unmarshal(body, &quotedata)
		
		quotedata[0].Content = strings.Replace(quotedata[0].Content, "<p>", "", -1)
		quotedata[0].Content = strings.Replace(quotedata[0].Content, "</p>", "", -1)
		return fmt.Sprintf("%s: %s", quotedata[0].Title, quotedata[0].Content)
	}
	return "what?"
}

func main() {
	args := os.Args
	if len(args) < 2 {
		panic("Give your api token")
	}
	botApiToken := args[1]

	bot, err := tgbotapi.NewBotAPI(botApiToken)

	if err != nil {
		panic(err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, message(update.Message.Text))
		msg.ReplyToMessageID = update.Message.MessageID

		bot.Send(msg)
	}
}
