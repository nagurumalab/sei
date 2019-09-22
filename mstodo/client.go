package mstodo

import (
	"context"
	"net/http"
	"time"

	"golang.org/x/oauth2"
)

//NewClient will init and give you back a todo client to work with
func NewClient() *http.Client {
	ctx := context.Background()
	httpClient := &http.Client{Timeout: 2 * time.Second}
	ctx = context.WithValue(ctx, oauth2.HTTPClient, httpClient)
	//TODO: Error handling code
	config := getOauthConfig()
	client, _ := userOAuthClient(ctx, config)
	return client
}
