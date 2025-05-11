package oidc

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
)

// Standard claims for OpenID Connect ID Tokens.
// https://openid.net/specs/openid-connect-core-1_0.html#IDToken
const (
	// Issuer Identifier for the Issuer of the response.
	Iss = "iss"

	// Subject Identifier: A locally unique and never reassigned identifier within the Issuer for the End-User.
	Sub = "sub"

	// Audience(s) that this ID Token is intended for.
	Aud = "aud"

	// Expiration time on or after which the ID Token MUST NOT be accepted.
	Exp = "exp"

	// Time at which the JWT was issued.
	Iat = "iat"

	// Time when the End-User authentication occurred.
	AuthTime = "auth_time"

	// String value used to associate a Client session with an ID Token and mitigate replay attacks.
	Nonce = "nonce"

	// Authentication Context Class Reference: Specifies the authentication context satisfied.
	Acr = "acr"

	// Authentication Methods References: Identifiers for authentication methods used.
	Amr = "amr"

	// Authorized party: The party to which the ID Token was issued.
	Azp = "azp"
)

func extractIDTokenPayload(idToken string) (map[string]any, error) {
	parts := strings.Split(idToken, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid ID token format")
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, fmt.Errorf("failed to decode payload: %w", err)
	}
	var claims map[string]any
	if err := json.Unmarshal(payload, &claims); err != nil {
		return nil, fmt.Errorf("failed to unmarshal payload: %w", err)
	}
	return claims, nil
}
