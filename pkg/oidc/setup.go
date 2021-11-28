package oidc

import (
	"crypto/tls"
	"net/http"

	"github.com/sirupsen/logrus"
)

// Setup comment lint rebel
func Setup(clientID, issuer, redirectURI string) (a App, err error) {
	a.Client = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	err = a.NewProvider(issuer)
	if err != nil {
		logrus.Error(err)
		return
	}

	var s struct {
		// What scopes does a provider support?
		//
		// See: https://openid.net/specs/openid-connect-discovery-1_0.html#ProviderMetadata
		ScopesSupported []string `json:"scopes_supported"`
	}
	if err = a.Provider.Claims(&s); err != nil {
		logrus.Errorf("failed to parse provider scopes_supported: %v", err)
		return
	}
	a.OfflineAsScope = true

	a.ClientID = clientID
	a.RedirectURI = redirectURI
	return
}
