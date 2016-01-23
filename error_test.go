package main

import (
	"testing"
)

func TestSlacksocError(t *testing.T) {
	var tests = []struct{
		given string
		expected string
	}{
		{"hi", "hi"},
	}

	for _, test := range tests {
		err := &SlacksocError{test.given}
		message := err.Error()
		if message != test.expected {
			t.Errorf("Failure. Expected: %s; Got: %s", test.expected, message)
		}
	}
}
