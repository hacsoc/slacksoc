package main

import (
	"net/url"
)

func (bot *Bot) GetChannelInfo() {
	if bot.Channels == nil {
		bot.Channels = make(map[string]string)
	}
	data := url.Values{"exclude_archived": []string{"1"}}
	resp, err := bot.Call("channels.list", data)
	payload, err := httpToJSON(resp, err)
	if err != nil {
		return
	}
	channels, ok := payload["channels"].([]interface{})
	if !ok {
		return
	}
	for _, channel := range channels {
		channelMap := channel.(JSONObject)
		id := channelMap["id"].(string)
		name := channelMap["name"].(string)
		bot.Channels[name] = id
	}
}
