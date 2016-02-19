package main

import (
	"fmt"
	"net/url"
	"os"

	"github.com/ajm188/slack"
	"github.com/ajm188/slack/plugins/github"
)

const (
	version = "0.4.0"
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
	payload := (<-userChan).(map[string]interface{})
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
	dm := <-dmChan
	return slack.NewMessage(text, dm), slack.Continue
}

func sendDM(bot *slack.Bot, event map[string]interface{}) (*slack.Message, slack.Status) {
	user := event["user"].(string)
	return bot.DirectMessage(user, "hi"), slack.Continue
}

func troll(bot *slack.Bot, event map[string]interface{}) (*slack.Message, slack.Status) {
	user, ok := event["user"]
	if !ok || user.(string) != bot.Users["catofnostalgia"] {
		return nil, slack.Continue
	}

	return bot.Mention(user.(string), "where is the third lambda?",
		event["channel"].(string)), slack.Continue
}

func configureGithubPlugin(id, secret, token string) {
	github.ClientID = id
	github.ClientSecret = secret
	github.AccessToken = token

	github.SharedClient = github.DefaultClient()
}

func getEnvvar(name string) (envvar string) {
	envvar = os.Getenv(name)
	if envvar == "" {
		fmt.Println("Missing environment variable %s", name)
		os.Exit(1)
	}
	return
}

func main() {
	token := getEnvvar("SLACKSOC_TOKEN")

	ghClientID := getEnvvar("GH_CLIENT_ID")
	ghClientSecret := getEnvvar("GH_CLIENT_SECRET")
	ghAccessToken := getEnvvar("GH_ACCESS_TOKEN")

	bot := slack.NewBot(token)
	bot.Respond("hi\\z", slack.Respond("hi there!"))
	bot.Respond("pm me", sendDM)
	bot.Respond("((what's)|(tell me) your)? ?version??",
		slack.Respond(fmt.Sprintf("My version is %s. My lib version is %s", version, slack.Version)))
	bot.Listen("gentoo", slack.React("funroll-loops"))
	bot.Listen(".+\\bslacksoc\\b", slack.React("raisedeyebrow"))
	bot.Listen("GNU/Linux", slack.React("stallman"))
	bot.OnEvent("message", troll)
	bot.OnEventWithSubtype("message", "channel_join", setRealNameFields)

	configureGithubPlugin(ghClientID, ghClientSecret, ghAccessToken)
	github.OpenIssue(bot, nil)

	fmt.Println("Starting bot")
	if err := bot.Start(); err != nil {
		fmt.Println(err)
	}
}
