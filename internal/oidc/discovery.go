package oidc

import (
	"encoding/json"
	"log"
	"sync"
	"time"
)

// DiscoveryData is referred to as "OpenID Provider Metadata" in the official
// specification: https://openid.net/specs/openid-connect-discovery-1_0.html#ProviderMetadata
type DiscoveryData struct {
	Issuer                                     string   `json:"issuer"`
	AuthorizationEndpoint                      string   `json:"authorization_endpoint"`
	TokenEndpoint                              string   `json:"token_endpoint,omitempty"`
	UserInfoEndpoint                           string   `json:"userinfo_endpoint,omitempty"`
	JwksURI                                    string   `json:"jwks_uri"`
	RegistrationEndpoint                       string   `json:"registration_endpoint,omitempty"`
	ScopesSupported                            []string `json:"scopes_supported,omitempty"`
	ResponseTypesSupported                     []string `json:"response_types_supported"`
	ResponseModesSupported                     []string `json:"response_modes_supported,omitempty"`
	GrantTypesSupported                        []string `json:"grant_types_supported,omitempty"`
	AcrValuesSupported                         []string `json:"acr_values_supported,omitempty"`
	SubjectTypesSupported                      []string `json:"subject_types_supported"`
	IDTokenSigningAlgValuesSupported           []string `json:"id_token_signing_alg_values_supported"`
	IDTokenEncryptionAlgValuesSupported        []string `json:"id_token_encryption_alg_values_supported,omitempty"`
	IDTokenEncryptionEncValuesSupported        []string `json:"id_token_encryption_enc_values_supported,omitempty"`
	UserInfoSigningAlgValuesSupported          []string `json:"userinfo_signing_alg_values_supported,omitempty"`
	UserInfoEncryptionAlgValuesSupported       []string `json:"userinfo_encryption_alg_values_supported,omitempty"`
	UserInfoEncryptionEncValuesSupported       []string `json:"userinfo_encryption_enc_values_supported,omitempty"`
	RequestObjectSigningAlgValuesSupported     []string `json:"request_object_signing_alg_values_supported,omitempty"`
	RequestObjectEncryptionAlgValuesSupported  []string `json:"request_object_encryption_alg_values_supported,omitempty"`
	RequestObjectEncryptionEncValuesSupported  []string `json:"request_object_encryption_enc_values_supported,omitempty"`
	TokenEndpointAuthMethodsSupported          []string `json:"token_endpoint_auth_methods_supported,omitempty"`
	TokenEndpointAuthSigningAlgValuesSupported []string `json:"token_endpoint_auth_signing_alg_values_supported,omitempty"`
	DisplayValuesSupported                     []string `json:"display_values_supported,omitempty"`
	ClaimTypesSupported                        []string `json:"claim_types_supported,omitempty"`
	ClaimsSupported                            []string `json:"claims_supported,omitempty"`
	ServiceDocumentation                       string   `json:"service_documentation,omitempty"`
	ClaimsLocalesSupported                     []string `json:"claims_locales_supported,omitempty"`
	UILocalesSupported                         []string `json:"ui_locales_supported,omitempty"`
	ClaimsParameterSupported                   bool     `json:"claims_parameter_supported,omitempty"`
	RequestParameterSupported                  bool     `json:"request_parameter_supported,omitempty"`
	RequestURISupported                        bool     `json:"request_uri_parameter_supported,omitempty"`
	RequireRequestURIRegistration              bool     `json:"require_request_uri_registration,omitempty"`
	OPPolicyURI                                string   `json:"op_policy_uri,omitempty"`
	OPTosURI                                   string   `json:"op_tos_uri,omitempty"`
}

type cachedDiscoveryData struct {
	validUntil *time.Time
	data       *DiscoveryData
}

const httpHeaderExpiresFormat = "Mon, 02 Jan 2006 15:04:05 MST"

var discoveryCache = map[string]*cachedDiscoveryData{}

var mu = &sync.Mutex{}
var discoveryCacheMutexes = map[string]*sync.Mutex{}

// GetDiscoveryData fetches the discovery data from the given URL and caches it.
func GetDiscoveryData(url string) (*DiscoveryData, error) {
	// First we'll check the cache without locking.
	if data, ok := discoveryCache[url]; ok {
		if data.validUntil == nil || time.Now().Before(*data.validUntil) {
			return data.data, nil
		}
	}

	// Wasn't in the cache, safely get a mutex for this url
	mu.Lock()
	urlSpecificMu := discoveryCacheMutexes[url]
	if urlSpecificMu == nil {
		urlSpecificMu = &sync.Mutex{}
		discoveryCacheMutexes[url] = urlSpecificMu
	}
	mu.Unlock()

	urlSpecificMu.Lock()
	defer urlSpecificMu.Unlock()

	// Check again after acquiring the lock
	if data, ok := discoveryCache[url]; ok {
		if data.validUntil == nil || time.Now().Before(*data.validUntil) {
			return data.data, nil
		}
	}

	// Couldn't find it in the cache, so we need to fetch it

	discoveryDataResp, err := HTTPClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer discoveryDataResp.Body.Close()

	expiresHeader := discoveryDataResp.Header.Get("Expires")

	var validUntil time.Time
	var nilifiedValidUntil *time.Time

	if expiresHeader != "" {
		validUntil, err = time.Parse(httpHeaderExpiresFormat, expiresHeader)
		if err != nil {
			log.Println("could not parse Expires header:", err)
		} else {
			nilifiedValidUntil = &validUntil
		}
	}

	var discoveryData DiscoveryData
	err = json.NewDecoder(discoveryDataResp.Body).Decode(&discoveryData)
	if err != nil {
		return nil, err
	}

	discoveryCache[url] = &cachedDiscoveryData{
		validUntil: nilifiedValidUntil,
		data:       &discoveryData,
	}

	return &discoveryData, nil
}
