package oidc

import (
	"context"
	"log"
	"os"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

var (
	provider     *oidc.Provider
	oauth2Config *oauth2.Config
)

// OIDC Provider 초기화
func InitOIDCProvider() error {
	ctx := context.Background()

	// Logto OIDC Provider 설정
	logtoEndpoint := os.Getenv("LOGTO_ENDPOINT")
	if logtoEndpoint == "" {
		return nil // OIDC 설정이 없으면 스킵
	}

	log.Printf("[OIDC] Initializing provider: %s", logtoEndpoint)

	var err error
	provider, err = oidc.NewProvider(ctx, logtoEndpoint)
	if err != nil {
		return err
	}

	// OAuth2 Config 설정
	oauth2Config = &oauth2.Config{
		ClientID:     os.Getenv("LOGTO_APP_ID"),
		ClientSecret: os.Getenv("LOGTO_APP_SECRET"),
		RedirectURL:  os.Getenv("LOGTO_REDIRECT_URI"),
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email", "roles"},
	}

	log.Println("[OIDC] Provider initialized successfully")
	return nil
}

// OAuth2 Config 반환
func GetOAuth2Config() *oauth2.Config {
	return oauth2Config
}

// OIDC Provider 반환
func GetProvider() *oidc.Provider {
	return provider
}

// OIDC가 설정되어 있는지 확인
func IsOIDCEnabled() bool {
	return provider != nil && oauth2Config != nil
}
