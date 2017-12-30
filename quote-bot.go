/**
 *	It's a simple telegram bot.
 * just a simple code to show how to work with telegram bots.
 * I hope comments help you well.
 * By Alex.
 */

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

// struct to use it for json we get from quote api.
type Quote struct {
	ID int 
	Title string
	Content string
	Link string 
}


func message(input string) string {
	// Show help message.
	if input == "help" || input == "/help" {
		return "There is some command: \n\n /quote -> Show a quote."
	}

	// Show welcome message
	if input == "start" || input == "/start" {
		return "Welcome to our bot!"
	}

	// Get a quote from api and return it.
	if input == "/quote" || input == "quote" {
		// Call the api and check for errors.
		resp, err := http.Get(QUOTE_URL)
		if err != nil {
			panic(err)
		}

		// Close the response body.
		defer resp.Body.Close()
		
		// Read the response body we got from quote api.
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		
		// unmarshal our json and pass it to a slice.
		// First element of slice is the quote we use.
		var quotedata  []Quote
		json.Unmarshal(body, &quotedata)
		
		// Delete the <p> tags from the quote.
		quotedata[0].Content = strings.Replace(quotedata[0].Content, "<p>", "", -1)
		quotedata[0].Content = strings.Replace(quotedata[0].Content, "</p>", "", -1)
		// Return the quote.
		return fmt.Sprintf("%s: %s", quotedata[0].Title, quotedata[0].Content)
	}
	// Non matched!!
	return "what?"
}

func main() {
	// Get the arguments so we can get the telegram bot's api token.
	args := os.Args
	if len(args) < 2 {
		panic("Give your api token")
	}
	botApiToken := args[1]

	// Create new bot
	bot, err := tgbotapi.NewBotAPI(botApiToken)

	// Check if there is any error.
	if err != nil {
		panic(err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	// Get updates on the bot.
	// It will return new messages.
	updates, err := bot.GetUpdatesChan(u)

	// Create a loop to review requests.
	for update := range updates {
		// continue if message is empty
		if update.Message == nil {
			continue
		}

		// Create a message to return it to the user.
		// `message` function will process the message and return the response message.
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, message(update.Message.Text))
		msg.ReplyToMessageID = update.Message.MessageID

		bot.Send(msg)
	}
}
