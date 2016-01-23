package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type JSONObject map[string]interface{}

func httpToJSON(resp *http.Response, err error) (JSONObject, error) {
	if err != nil {
		return nil, err
	}
	raw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var payload JSONObject
	err = json.Unmarshal(raw, &payload)
	if err != nil {
		return nil, err
	}
	return payload, nil
}
