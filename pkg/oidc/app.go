package oidc

import (
	"context"
	"net/http"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

// App comment lint rebel
type App struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string

	Verifier *oidc.IDTokenVerifier
	Provider *oidc.Provider

	// Does the provider use "offline_access" scope to request a refresh token
	// or does it use "access_type=offline" (e.g. Google)?
	OfflineAsScope bool

	Client *http.Client
}

// Oauth2Config comment lint rebel
func (a *App) Oauth2Config(scopes []string) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     a.ClientID,
		ClientSecret: a.ClientSecret,
		Endpoint:     a.Provider.Endpoint(),
		Scopes:       scopes,
		RedirectURL:  a.RedirectURI,
	}
}

const exampleAppState = "ianvs"

// GetLoginURL comment lint rebel
func (a *App) GetLoginURL() string {
	var scopes []string

	authCodeURL := ""
	scopes = append(scopes, "openid", "profile", "email")
	if a.OfflineAsScope {
		scopes = append(scopes, "offline_access")
		authCodeURL = a.Oauth2Config(scopes).AuthCodeURL(exampleAppState)
	} else {
		authCodeURL = a.Oauth2Config(scopes).AuthCodeURL(exampleAppState, oauth2.AccessTypeOffline)
	}

	return authCodeURL
}

// GetContext comment lint rebel
func (a *App) GetContext() context.Context {
	return oidc.ClientContext(context.Background(), a.Client)
}

// NewProvider comment lint rebel
func (a *App) NewProvider(issuer string) (err error) {
	a.Provider, err = oidc.NewProvider(a.GetContext(), issuer)
	return
}

// GetVerifier comment lint rebel
func (a *App) GetVerifier() *oidc.IDTokenVerifier {
	if a.Verifier == nil {
		a.Verifier = a.Provider.Verifier(&oidc.Config{ClientID: a.ClientID})
	}
	return a.Verifier
}

// GetToken comment lint rebel
func (a *App) GetToken(refreshToken string) (token *oauth2.Token, err error) {
	oauth2Config := a.Oauth2Config(nil)
	t := &oauth2.Token{
		RefreshToken: refreshToken,
		Expiry:       time.Now().Add(-time.Hour),
	}
	token, err = oauth2Config.TokenSource(a.GetContext(), t).Token()
	return
}

// Exchange comment lint rebel
func (a *App) Exchange(code string) (token *oauth2.Token, err error) {
	oauth2Config := a.Oauth2Config(nil)
	token, err = oauth2Config.Exchange(a.GetContext(), code)
	return
}

// FetchCode comment lint rebel
func (a *App) FetchCode() {

}
