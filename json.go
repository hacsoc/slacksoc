package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func httpToJSON(resp *http.Response, err error) (map[string]interface{}, error) {
	if err != nil {
		return nil, err
	}
	raw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var payload map[string]interface{}
	err = json.Unmarshal(raw, &payload)
	return payload, err
}
