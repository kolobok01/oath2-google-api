package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// oauthClient shows how to use an OAuth client ID to authenticate as an end-user.
func oauthClient() error {
	ctx := context.Background()

	redirectURL := "urn:ietf:wg:oauth:2.0:oob"
	config := &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  redirectURL,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}

	// Dummy authorization flow to read auth code from stdin.
	authURL := config.AuthCodeURL("state", oauth2.AccessTypeOffline)
	fmt.Printf("Follow the link in your browser to obtain auth code (AuthURL): %s\n", authURL)

	fmt.Printf("Input auth code:")
	// Read the authentication code from the command line
	var code string
	fmt.Scanln(&code)

	// Exchange auth code for OAuth token.
	token, err := config.Exchange(ctx, code)
	if err != nil {
		fmt.Println("Exchange Error, err = ", err)
		return fmt.Errorf("config.Exchange: %v", err)
	}

	fmt.Println("token.AccessToken:", token.AccessToken)
	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		fmt.Println("http Get Error, err = ", err)
		return err
	}
	defer response.Body.Close()
	// fmt.Printf("[getUserInfo] response.Body=%v\n", response.Body)
	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Read Error, err = ", err)
	}
	bodyString := string(bodyBytes)
	fmt.Println(bodyString)

	return nil
}

func main() {
	oauthClient()
}
