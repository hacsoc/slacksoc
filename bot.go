package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

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
	Channels []string
	WebSocketURL string
}

func NewBot(token string, channels []string) *Bot {
	return &Bot{Token: token, Channels: channels}
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
	}
}

func main() {
	token := os.Getenv(TOKEN_VAR)
	if token == "" {
		fmt.Println(NO_TOKEN_ERROR)
		os.Exit(1)
	}

	bot := NewBot(token, nil)
	fmt.Println("Starting bot")
	if err := bot.Start(); err != nil {
		fmt.Println(err)
	}
	fmt.Println("Looping")
	if err := bot.Loop(); err != nil {
		fmt.Println(err)
	}
}
