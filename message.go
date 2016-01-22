package main

import (
	"time"
)

func mentionText(nick, beforeNick, afterNick string) string {
	text := beforeNick
	nick = "<@" + nick + ">"
	text += nick
	if text == nick && afterNick != "" {
		text += ":"
	}
	if afterNick != "" {
		text += " " + afterNick
	}
	return text
}

func Mention(nick, channel, beforeNick, afterNick string) interface{} {
	text := mentionText(nick, beforeNick, afterNick)
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
