package main

import (
	"time"
)

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
