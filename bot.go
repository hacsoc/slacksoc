package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/hacsoc/slacksoc/api"
)

const (
	TOKEN_VAR = "SLACKSOC_TOKEN"
	NO_TOKEN_ERROR = "You must have the SLACKSOC_TOKEN variable to run the" +
					 " slacksoc bot"
)

type Bot struct {
	Token string
	Channels map[string]string
	WebSocketURL string
}

func NewBot(token string) *Bot {
	return &Bot{Token: token}
}

func (bot *Bot) Call(method string, data url.Values) (*http.Response, error) {
	data.Set("token", bot.Token)
	return api.Call(method, data)
}

func (bot *Bot) Start() error {
	payload, err := httpToJSON(bot.Call("rtm.start", url.Values{}))
	if err != nil {
		return err
	}
	ok, present := payload["ok"].(bool)
	if !present || ok != true {
		return &RTMStartError{"could not connect to RTM API"}
	}
	bot.GetChannelInfo()
	websocketURL, _ := payload["url"].(string)
	bot.WebSocketURL = websocketURL
	return nil
}

func (bot *Bot) Loop() error {
	dialer := websocket.Dialer{}
	conn, _, err := dialer.Dial(bot.WebSocketURL, http.Header{})
	if err != nil {
		return err
	}
	for {
		messageType, bytes, err := conn.ReadMessage()
		if err != nil {
			// NextReader returns an error if the connection is closed
			conn.Close()
			return nil
		}
		if messageType == websocket.BinaryMessage {
			continue // ignore binary messages
		}
		var message map[string]interface{}
		if err = json.Unmarshal(bytes, &message); err != nil {
			continue
		}
		fmt.Println("", message)
		if _, ok := message["type"]; !ok {
			continue
		}
		switch message["type"].(string) {
		case "message":
			bot.ReceiveMessage(conn, message)
		default:
			continue
		}
	}
}

func (bot *Bot) ReceiveMessage(conn *websocket.Conn, message map[string]interface{}) {
	subtype, _ := message["subtype"]
	hiddenSubtype, ok := message["hidden"]
	hidden := ok && hiddenSubtype.(bool)
	reply := bot.ConstructReply(message, subtype, hidden)
	if reply != nil {
		conn.WriteJSON(reply)
	}
}

func (bot *Bot) ConstructReply(message map[string]interface{}, subtype interface{}, hidden bool) interface{} {
	if subtype != nil {
		switch subtype.(string) {
		case "bot_message":
			// don't reply to other bots
			return nil
		case "channel_join":
			return bot.SetRealNameFields(message)
		default:
			return nil
		}
	} else {
		text := message["text"].(string)
		if strings.Contains(text, "hi slacksoc") {
			return Mention(message["user"].(string), message["channel"].(string), "hi ", "")
		} else if text == "slacksoc: pm me" {
			return bot.DirectMessage(message["user"].(string), "hi")
		}
		return nil
	}
}

func (bot *Bot) SetRealNameFields(message map[string]interface{}) interface{} {
	channel := message["channel"].(string)
	if channel != bot.Channels["general"] {
		return nil
	}
	userID := message["user"].(string)
	dmChan := make(chan string)
	userChan := make(chan interface{})
	go func() {
		dm, _ := bot.OpenDirectMessage(userID)
		dmChan <- dm
	}()
	go func() {
		resp, err := bot.Call("users.info", url.Values{"user": []string{userID}})
		payload, err := httpToJSON(resp, err)
		userChan <- payload
	}()
	payload := (<- userChan).(map[string]interface{})
	success := payload["ok"].(bool)
	if !success {
		fmt.Println(payload)
		return nil
	}
	user := payload["user"].(map[string]interface{})
	nick := user["name"].(string)
	text := "Please set your real name fields. https://hacsoc.slack.com/team/%s."
	text += " Then click \"Edit\"."
	text = fmt.Sprintf(text, nick)
	dm := <- dmChan
	return Message(text, dm)
}

func (bot *Bot) DirectMessage(user, text string) interface{} {
	dm, err := bot.OpenDirectMessage(user)
	if err != nil {
		return nil
	}
	return Message(text, dm)
}

func (bot *Bot) OpenDirectMessage(user string) (string, error) {
	resp, err := bot.Call("im.open", url.Values{"user": []string{user}})
	payload, err := httpToJSON(resp, err)
	if err != nil {
		return "", err
	}
	success := payload["ok"].(bool)
	if !success {
		fmt.Println(payload)
		return "", nil // need an actual error here
	}
	channel := payload["channel"].(map[string]interface{})
	return channel["id"].(string), nil
}

func Mention(nick, channel, beforeNick, afterNick string) interface{} {
	text := beforeNick
	nick = "<@" + nick + ">"
	text += nick
	if text == nick {
		text += ": "
	}
	text += afterNick
	return Message(text, channel)
}

func Message(text, channel string) interface{} {
	return map[string]string{
		"id": time.Now().Format("010206150405"),
		"type": "message",
		"channel": channel,
		"text": text,
	}
}

func main() {
	token := os.Getenv(TOKEN_VAR)
	if token == "" {
		fmt.Println(NO_TOKEN_ERROR)
		os.Exit(1)
	}

	bot := NewBot(token)
	fmt.Println("Starting bot")
	if err := bot.Start(); err != nil {
		fmt.Println(err)
	}
	fmt.Println("Looping")
	if err := bot.Loop(); err != nil {
		fmt.Println(err)
	}
}
