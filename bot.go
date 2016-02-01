package main

import (
	"fmt"
	"net/url"
	"os"

	"github.com/ajm188/slack"
)

const (
	tokenVar = "SLACKSOC_TOKEN"
	noTokenError = "You must have the SLACKSOC_TOKEN variable to run the" +
					 " slacksoc bot"
	version = "0.1.0"
)

func setRealNameFields(bot *slack.Bot, event map[string]interface{}) (*slack.Message, slack.Status) {
	channel := event["channel"].(string)
	if channel != bot.Channels["general"] {
		return nil, slack.Continue
	}
	userID := event["user"].(string)
	dmChan := make(chan string)
	userChan := make(chan interface{})
	go func() {
		dm, _ := bot.OpenDirectMessage(userID)
		dmChan <- dm
	}()
	go func() {
		payload, _ := bot.Call("users.info", url.Values{"user": []string{userID}})
		userChan <- payload
	}()
	payload := (<- userChan).(map[string]interface{})
	success := payload["ok"].(bool)
	if !success {
		fmt.Println(payload)
		return nil, slack.Continue
	}
	user := payload["user"].(map[string]interface{})
	nick := user["name"].(string)
	text := "Please set your real name fields. https://hacsoc.slack.com/team/%s."
	text += " Then click \"Edit\"."
	text = fmt.Sprintf(text, nick)
	dm := <- dmChan
	return slack.NewMessage(text, dm), slack.Continue
}

func sendDM(bot *slack.Bot, event map[string]interface{}) (*slack.Message, slack.Status) {
	user := event["user"].(string)
	return bot.DirectMessage(user, "hi"), slack.Continue
}

func main() {
	token := os.Getenv(tokenVar)
	if token == "" {
		fmt.Println(noTokenError)
		os.Exit(1)
	}

	bot := slack.NewBot(token)
	bot.Respond("hi", slack.Respond("hi there!"))
	bot.Respond("pm me", sendDM)
	bot.Respond("((what's)|(tell me) your)? ?version??",
		slack.Respond(fmt.Sprintf("My version is %s. My lib version is %s", version, slack.Version)))
	bot.Listen("gentoo", slack.React("funroll-loops"))
	bot.OnEventWithSubtype("message", "channel_join", setRealNameFields)
	fmt.Println("Starting bot")
	if err := bot.Start(); err != nil {
		fmt.Println(err)
	}
}
