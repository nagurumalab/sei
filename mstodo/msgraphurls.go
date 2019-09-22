package mstodo

import (
	"net/url"
	"path"
)

//BASEURL is the base of all urls -_-
var BASEURL = "https://graph.microsoft.com/beta/me/outlook/"

func constructURL(pathSegments []string, params map[string]string) string {
	reqURL, _ := url.Parse(BASEURL)
	reqURL.Path += path.Join(pathSegments...)
	values := url.Values{}
	for key, value := range params {
		values.Add(key, value)
	}
	reqURL.RawQuery = values.Encode()
	return reqURL.String()
}
