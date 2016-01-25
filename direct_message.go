package main

import (
	"fmt"
	"net/url"
)

func (bot *Bot) DirectMessage(user, text string) interface{} {
	dm, err := bot.OpenDirectMessage(user)
	if err != nil {
		return nil
	}
	return NewMessage(text, dm).ToMap()
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
	channel := payload["channel"].(JSONObject)
	return channel["id"].(string), nil
}
