package oidc

// federatedIDClaims comment lint rebel
type federatedIDClaims struct {
	ConnectorID string `json:"connector_id,omitempty"`
	UserID      string `json:"user_id,omitempty"`
}

// Claims comment lint rebel
type Claims struct {
	Issuer           string `json:"iss"`
	Subject          string `json:"sub"`
	Expiry           int64  `json:"exp"`
	IssuedAt         int64  `json:"iat"`
	AuthorizingParty string `json:"azp,omitempty"`
	Nonce            string `json:"nonce,omitempty"`

	AccessTokenHash string `json:"at_hash,omitempty"`
	CodeHash        string `json:"c_hash,omitempty"`

	Email         string `json:"email,omitempty"`
	EmailVerified *bool  `json:"email_verified,omitempty"`

	Groups []string `json:"groups,omitempty"`

	Name              string `json:"name,omitempty"`
	PreferredUsername string `json:"preferred_username,omitempty"`

	FederatedIDClaims *federatedIDClaims `json:"federated_claims,omitempty"`
}
