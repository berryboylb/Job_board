package auth

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	mrand "math/rand"
	"net/http"
	"os"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-contrib/sessions"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"gorm.io/gorm"

	"job_board/db"
	"job_board/models"
)

var database *gorm.DB

func init() {
	database = db.GetDB()
}

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

func generateToken(sub string) (*TokenResponse, error) {
	url := "https://" + os.Getenv("AUTH0_DOMAIN") + "/oauth/token"

	clientID := os.Getenv("TOKEN_CLIENT_ID")
	clientSecret := os.Getenv("TOKEN_CLIENT_SECRET")
	audience := os.Getenv("TOKEN_AUDIENCE")

	if clientID == "" || clientSecret == "" || audience == "" {
		return nil, fmt.Errorf("environment variable missing")
	}

	data := map[string]interface{}{
		"client_id":     clientID,
		"client_secret": clientSecret,
		"audience":      audience,
		"grant_type":    "client_credentials",
		"sub":           sub,
	}

	payload, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
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

func GoogleUser(session sessions.Session) (*models.User, bool, error) {
	profile, ok := session.Get("profile").(GoogleResponse)
	if !ok {
		return nil, false, fmt.Errorf("profile data not found or not  valid for google response")
	}

	userType, ok := session.Get("type").(models.RoleAllowed)
	if !ok {
		return nil, false, fmt.Errorf("invalid type")
	}
	subscriberId := uuid.NewString()
	user := models.User{
		Email:        profile.Email,
		Name:         profile.Name,
		Picture:      profile.Picture,
		ProviderID:   profile.Sub,
		RoleName:     userType,
		SubscriberID: subscriberId,
	}
	return CreateUser(user)
}

func EmailUser(session sessions.Session) (*models.User, bool, error) {
	profile, ok := session.Get("profile").(EmailResponse)
	if !ok {
		return nil, false, fmt.Errorf("profile data not found or not  valid for email response")
	}

	userType, ok := session.Get("type").(models.RoleAllowed)
	if !ok {
		return nil, false, fmt.Errorf("invalid type")
	}
	subscriberId := uuid.NewString()
	user := models.User{
		Email:        profile.Email,
		Name:         profile.Name,
		Picture:      profile.Picture,
		ProviderID:   profile.Sub,
		RoleName:     userType,
		SubscriberID: subscriberId,
	}
	return CreateUser(user)
}

func GithubUser(session sessions.Session) (*models.User, bool, error) {
	profile, ok := session.Get("profile").(GithubResponse)
	if !ok {
		return nil, false, fmt.Errorf("profile data not found or not  valid for email response")
	}
	userType, ok := session.Get("type").(models.RoleAllowed)
	if !ok {
		return nil, false, fmt.Errorf("invalid type")
	}
	subscriberId := uuid.NewString()

	user := models.User{
		Name:         profile.Name,
		Picture:      profile.Picture,
		ProviderID:   profile.Sub,
		RoleName:     userType,
		SubscriberID: subscriberId,
	}
	return CreateUser(user)
}

func handleUser(sub string, session sessions.Session) (*models.User, bool, error) {
	switch sub {
	case "google-oauth2":
		return GoogleUser(session)
	case "auth0":
		return EmailUser(session)
	case "github":
		return GithubUser(session)
	default:
		return nil, false, fmt.Errorf("unsupported OAuth provider: %s", sub)
	}
}

func CreateUser(user models.User) (*models.User, bool, error) {
	// Start a new transaction
	tx := database.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Try to find the user
	existingUser := &models.User{}
	tx = tx.Where("provider_id = ?", user.ProviderID)
	if user.Email != "" {
		tx = tx.Or("email = ?", user.Email)
	}
	result := tx.Unscoped().Preload("Profile").First(existingUser)
	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			tx.Rollback()
			return nil, false, fmt.Errorf("error fetching user: %w", result.Error)
		}
	}

	if result.RowsAffected > 0 {
		// User already exists;  return them
		if err := tx.Commit().Error; err != nil {
			return nil, false, fmt.Errorf("error committing transaction: %w", err)
		}
		fmt.Println("THE TWO IDS ARE THE SAME",  existingUser.ProviderID, user.ProviderID )
		if existingUser.ProviderID != user.ProviderID {
			return nil, false, fmt.Errorf("error returning user: Provider ID mismatch for known user, please use a different provider or a different email %v", nil)
		}
		if existingUser.DeletedAt.Valid {
			return nil, false, fmt.Errorf("error returning user: your account was deleted, but hasn't been cleared, contact support to reinstate it or it would be permanently cleared after 3 months %w", nil)
		}
		return existingUser, false, nil
	}

	// User not found, so create it
	if err := tx.WithContext(context.Background()).Create(&user).Error; err != nil {
		tx.Rollback()
		return nil, false, fmt.Errorf("error creating a new user: %v", err.Error())
	}
	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return nil, false, fmt.Errorf("error committing transaction: %w", err)
	}
	return &user, true, nil
}

const charset = "0123456789"

var seededRand *mrand.Rand = mrand.New(
	mrand.NewSource(time.Now().UnixNano()))

func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	fmt.Println(string(b), "otp")
	return string(b)
}

func GenerateOtp(length int) string {
	return StringWithCharset(length, charset)
}
