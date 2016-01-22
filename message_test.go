package main

import (
	"testing"
)

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
