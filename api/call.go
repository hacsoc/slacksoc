package api

import (
	"fmt"
	"net/http"
	"net/url"
)

const (
	SLACK_URL = "https://slack.com/api/%s"
)

func Call(method string, data url.Values) (*http.Response, error) {
	methodURL := fmt.Sprintf(SLACK_URL, method)
	return http.PostForm(methodURL, data)
}
