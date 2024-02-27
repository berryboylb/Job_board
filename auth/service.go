package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-contrib/sessions"
	"golang.org/x/oauth2"
)

// Authenticator is used to authenticate our users.
type Authenticator struct {
	*oidc.Provider
	oauth2.Config
}

// New instantiates the *Authenticator.
func New() (*Authenticator, error) {
	provider, err := oidc.NewProvider(
		context.Background(),
		"https://"+os.Getenv("AUTH0_DOMAIN")+"/",
	)
	if err != nil {
		return nil, err
	}

	conf := oauth2.Config{
		ClientID:     os.Getenv("AUTH0_CLIENT_ID"),
		ClientSecret: os.Getenv("AUTH0_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("AUTH0_CALLBACK_URL"),
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}

	return &Authenticator{
		Provider: provider,
		Config:   conf,
	}, nil
}

// VerifyIDToken verifies that an *oauth2.Token is a valid *oidc.IDToken.
func (a *Authenticator) VerifyIDToken(ctx context.Context, token *oauth2.Token) (*oidc.IDToken, error) {
	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		return nil, errors.New("no id_token field in oauth2 token")
	}

	oidcConfig := &oidc.Config{
		ClientID: a.ClientID,
	}

	return a.Verifier(oidcConfig).Verify(ctx, rawIDToken)
}

func generateRandomState() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	state := base64.StdEncoding.EncodeToString(b)

	return state, nil
}

func generateToken() (*TokenResponse, error) {
	url := "https://" + os.Getenv("AUTH0_DOMAIN") + "/oauth/token"

	clientID := os.Getenv("TOKEN_CLIENT_ID")
	clientSecret := os.Getenv("TOKEN_CLIENT_SECRET")
	audience := "https://" + os.Getenv("AUTH0_DOMAIN") + "/api/v2/"

	if clientID == "" || clientSecret == "" || audience == "" {
		return nil, fmt.Errorf("environment variable missing")
	}

	payload := strings.NewReader(fmt.Sprintf("{\"client_id\":\"%s\",\"client_secret\":\"%s\",\"audience\":\"%s\",\"grant_type\":\"client_credentials\"}", clientID, clientSecret, audience))

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("content-type", "application/json")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var response TokenResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func decoder(idToken *oidc.IDToken, sub string) (interface{}, error) {
	var profile interface{}
	switch sub {
	case "google-oauth2":
		profile = &GoogleResponse{}
	case "auth0":
		profile = &EmailResponse{}
	case "github":
		profile = &GithubResponse{}
	default:
		return nil, fmt.Errorf("unsupported OAuth provider: %s", sub)
	}

	if err := idToken.Claims(&profile); err != nil {
		return nil, err
	}
	return profile, nil
}

func GoogleUser(session sessions.Session) (interface{}, error) {
	profile, ok := session.Get("profile").(GoogleResponse)
	if !ok {
		return nil, fmt.Errorf("Profile data not found or not  valid for google response")
	}

	userType, ok := session.Get("type").(string)
	if !ok {
		return nil, fmt.Errorf("invalid type")
	}
	fmt.Println(userType)
	//create db user here
	return profile, nil
}

func EmailUser(session sessions.Session) (interface{}, error) {
	profile, ok := session.Get("profile").(EmailResponse)
	if !ok {
		return nil, fmt.Errorf("profile data not found or not  valid for email response")
	}

	userType, ok := session.Get("type").(string)
	if !ok {
		return nil, fmt.Errorf("invalid type")
	}
	
	fmt.Println(userType)
	//create db user here
	return profile, nil
}

func GithubUser(session sessions.Session) (interface{}, error) {
	profile, ok := session.Get("profile").(GithubResponse)
	if !ok {
		return nil, fmt.Errorf("profile data not found or not  valid for email response")
	}
	userType, ok := session.Get("type").(string)
	if !ok {
		return nil, fmt.Errorf("invalid type")
	}
	fmt.Println(userType)
	//create db user here
	return profile, nil
}

func handleUser(sub string, session sessions.Session) (interface{}, error) {
	switch sub {
	case "google-oauth2":
		return GoogleUser(session)
	case "auth0":
		return EmailUser(session)
	case "github":
		return GithubUser(session)
	default:
		return nil, fmt.Errorf("unsupported OAuth provider: %s", sub)
	}
}
