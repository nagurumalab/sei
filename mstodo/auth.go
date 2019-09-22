package mstodo

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/microsoft"
)

var tokenStore = "token.json"

//userOAuthClient return a client
func userOAuthClient(ctx context.Context, config *oauth2.Config) (client *http.Client, err error) {
	var userToken *oauth2.Token
	if userToken, err = getCachedToken(); err != nil {
		// if token for user is not cached then go through oauth2 flow
		if userToken, err = newUserToken(ctx, config); err != nil {
			log.Fatal(err)
			return
		}
	}
	if !userToken.Valid() { // if user token is expired
		userToken = &oauth2.Token{RefreshToken: userToken.RefreshToken}
	}
	return config.Client(ctx, userToken), err
}

func getCachedToken() (*oauth2.Token, error) {
	tk := new(oauth2.Token)
	if _, err := os.Stat(tokenStore); os.IsNotExist(err) {
		log.Print(err)
		return nil, err
	}
	jToken, err := ioutil.ReadFile(tokenStore)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	err = json.Unmarshal(jToken, tk)
	if err != nil {
		log.Panic(err)
	}
	return tk, err
}

func saveToken(token *oauth2.Token) error {
	jToken, err := json.Marshal(token)
	if err != nil {
		log.Fatal(err)
	}
	return ioutil.WriteFile(tokenStore, jToken, 0400)
}

func newUserToken(ctx context.Context, config *oauth2.Config) (*oauth2.Token, error) {
	stateBytes := make([]byte, 32)
	_, err := rand.Read(stateBytes)
	if err != nil {
		log.Fatalf("Unable to read random bytes: %v", err)
		return nil, err
	}
	state := fmt.Sprintf("%x", stateBytes)
	// Redirect user to consent page to ask for permission
	// for the scopes specified above.
	url := config.AuthCodeURL(state, oauth2.AccessTypeOffline)
	fmt.Printf("Visit the URL for the auth dialog: %v", url)

	// Use the authorization code that is pushed to the redirect
	// URL. Exchange will do the handshake to retrieve the
	// initial access token. The HTTP Client returned by
	// conf.Client will refresh the token as necessary.
	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatal(err)
	}

	// Use the custom HTTP client when requesting a token.
	token, err := config.Exchange(ctx, authCode)
	if err != nil {
		log.Fatalf("Exchange error: %v", err)
		return nil, err
	}

	saveToken(token) // save token to datastore

	return token, nil
}

func getOauthConfig() *oauth2.Config {

	// TODO: remove the client id  and secret from the code and put it in a config file
	conf := &oauth2.Config{
		ClientID:    "c7bcff77-9645-458a-a39e-f5d3c2664ea8",
		Scopes:      []string{"Tasks.Read", "Tasks.Read.Shared", "Tasks.ReadWrite", "Tasks.ReadWrite.Shared", "User.Read", "offline_access"},
		Endpoint:    microsoft.AzureADEndpoint("common"),
		RedirectURL: "https://login.microsoftonline.com/common/oauth2/nativeclient",
	}
	return conf
}
