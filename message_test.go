package main

import (
	"testing"
)

func TestMessageToMap(t *testing.T) {
	var tests = []struct {
		message *Message
		expected map[string]string
	}{
		{
			&Message{"id", "type", "channel", "text"},
		 	map[string]string{
				"id": "id",
				"type": "type",
				"channel": "channel",
				"text": "text",
			},
		},
	}

	for _, test := range tests {
		result := test.message.ToMap()
		if !mapsAreEqual(result, test.expected) {
			t.Error("Failure. Expecting: ", test.expected, "; Got: ", result)
		}
	}
}

func TestNewMessage(t *testing.T) {
	message := NewMessage("message text", "#general")
	if message.text != "message text" {
		t.Errorf("Failure. Expecting: \"message text\"; Got: %s", message.text)
	}
	if message.channel != "#general" {
		t.Errorf("Failure. Expecting: \"#general\"; Got: %s", message.channel)
	}
	if message.messageType != "message" {
		t.Errorf("Failure. Message did not have messageType \"message\"")
	}
}

func TestMention(t *testing.T) {
	var tests = []struct {
		before string
		after string
		expected string
	}{
		{"", "", "<@hi>"},
		{"", "hello", "<@hi>: hello"},
		{"hello ", "", "hello <@hi>"},
		{"hello ", "world", "hello <@hi> world"},
	}

	for _, test := range tests {
		text := mentionText("hi", test.before, test.after)
		if test.expected != text {
			t.Errorf("Failure. Expecting: %s; Got: %s", test.expected, text)
		}
	}
}

func mapsAreEqual(map1, map2 map[string]string) bool {
	if len(map1) != len(map2) {
		return false
	}
	for key, value1 := range map1 {
		value2, ok := map2[key]
		if !ok || value1 != value2 {
			return false
		}
	}
	return true
}
