package main

import (
	"time"
)

type Message struct {
	id string
	messageType string
	channel string
	text string
}

func (message *Message) ToMap() map[string]string {
	return map[string]string{
		"id": message.id,
		"type": message.messageType,
		"channel": message.channel,
		"text": message.text,
	}
}

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
	return NewMessage(text, channel).ToMap()
}

func NewMessage(text, channel string) *Message {
	return &Message{
		id: time.Now().Format("010206150405"),
		messageType: "message",
		channel: channel,
		text: text,
	}
}
