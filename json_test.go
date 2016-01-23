package main

import (
	"testing"
)

type testError struct{}

func (err *testError) Error() string {
	return "error"
}

func TestHttpToJSON(t *testing.T) {
	// httpToJSON should do nothing when handed an error
	test := &testError{}
	json, err := httpToJSON(nil, test)
	if json != nil {
		t.Error("Failure. Expecting: nil; Got: ", json)
	}
	if err != test {
		t.Error("Failure. Expecting: ", test, "; Got: ", err)
	}
}
