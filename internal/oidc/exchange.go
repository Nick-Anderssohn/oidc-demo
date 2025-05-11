package oidc

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/oauth2"
)

type Config struct {
	*oauth2.Config
	DiscoveryURL string
}

type TokenResponse struct {
	*oauth2.Token
	IDToken        string
	IDTokenPayload map[string]any
}

func ExchangeCodeForToken(ctx context.Context, config *Config, code string, opts ...oauth2.AuthCodeOption) (TokenResponse, error) {
	updatedCTX := context.WithValue(ctx, oauth2.HTTPClient, HTTPClient)

	token, err := config.Exchange(updatedCTX, code, opts...)
	if err != nil {
		return TokenResponse{}, err
	}

	if token.TokenType != "Bearer" {
		return TokenResponse{}, fmt.Errorf("unexpected token type: %s", token.TokenType)
	}

	idTokenStr, ok := token.Extra("id_token").(string)
	if !ok {
		return TokenResponse{}, err
	}

	// We don't need to validate the signature of the jwt. According to the OIDC
	// spec, TLS server validation with the identity provider is sufficient. See
	// https://openid.net/specs/openid-connect-core-1_0.html#IDTokenValidation
	// Since we don't need to worry about the signature, we'll just validate the
	// payload claims and call it good.
	claims, err := extractIDTokenPayload(idTokenStr)
	if err != nil {
		return TokenResponse{}, err
	}
	if err := validateIDTokenStandardPayloadClaims(config, claims); err != nil {
		return TokenResponse{}, err
	}

	return TokenResponse{
		Token:          token,
		IDToken:        idTokenStr,
		IDTokenPayload: claims,
	}, nil
}

// https://openid.net/specs/openid-connect-core-1_0.html#IDTokenValidation
func validateIDTokenStandardPayloadClaims(config *Config, claims map[string]any) error {
	discoveryData, err := GetDiscoveryData(config.DiscoveryURL)
	if err != nil {
		return err
	}

	// Validate Issuer
	if iss, ok := claims[Iss].(string); !ok || iss != discoveryData.Issuer {
		return fmt.Errorf("invalid issuer: %v", claims[Iss])
	}

	audAsString, audIsStr := claims[Aud].(string)
	if audIsStr && audAsString != config.ClientID {
		return fmt.Errorf("aud does not match client ID. aud: %v", claims[Aud])
	}

	if !audIsStr {
		aud, ok := claims[Aud].([]any)
		if !ok || len(aud) == 0 {
			return fmt.Errorf("invalid audience: %v", claims[Aud])
		}

		clientIDFound := false
		for _, a := range aud {
			if audStr, ok := a.(string); ok && audStr == config.ClientID {
				clientIDFound = true
				break
			}
		}
		if !clientIDFound {
			return fmt.Errorf("client ID not found in audience: %v", claims[Aud])
		}
	}

	if exp, ok := claims[Exp].(float64); !ok || int64(exp) < time.Now().Unix() {
		return fmt.Errorf("token expired")
	}

	return nil
}
