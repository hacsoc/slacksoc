package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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

func (bot *Bot) call(method string, data url.Values) (*http.Response, error) {
	data.Set("token", bot.Token)
	return api.Call(method, data)
}

func (bot *Bot) startRTM() ([]byte, error) {
	resp, err := bot.call("rtm.start", url.Values{})
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(resp.Body)
}

func (bot *Bot) Start() error {
	rawData, err := bot.startRTM()
	if err != nil {
		return err
	}
	var payload map[string]interface{}
	err = json.Unmarshal(rawData, &payload)
	if err != nil {
		return err
	}
	ok, present := payload["ok"].(bool)
	if !present || ok != true {
		return &RTMStartError{"could not connect to RTM API"}
	}
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
			ReceiveMessage(conn, message)
		default:
			continue
		}
	}
}

func ReceiveMessage(conn *websocket.Conn, message map[string]interface{}) {
	subtype, hasSubtype := message["subtype"]
	hiddenSubtype, ok := message["hidden"]
	hidden := ok && hiddenSubtype.(bool)
	reply := ConstructReply(message, subtype, hasSubtype, hidden)
	if reply != nil {
		conn.WriteJSON(reply)
	}
}

func ConstructReply(message map[string]interface{}, subtype interface{}, hasSubtype, hidden bool) interface{} {
	if hasSubtype {
		switch subtype.(string) {
		case "bot_message":
			// don't reply to other bots
			return nil
		case "channel_join":
		}
		return nil
	} else {
		text := message["text"].(string)
		if strings.Contains(text, "hi slacksoc") {
			return Mention(message, "hi ", "")
		}
		return nil
	}
}

func Mention(message map[string]interface{}, beforeNick, afterNick string) interface{} {
	text := beforeNick
	nick := "<@" + message["user"].(string) + ">"
	text += nick
	if text == nick {
		text += ": "
	}
	text += afterNick
	j := map[string]string{
		"id": time.Now().Format("010206150405"),
		"type": "message",
		"channel": message["channel"].(string),
		"text": text,
	}
	return j
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
