package oidc

import (
	"net/http"
	"time"
)

// HttpClient is the http client used by this package. By default,
// it has some connection pooling. This client can be replaced with
// a custom one if needed.
var HttpClient = http.Client{
	Transport: &http.Transport{
		MaxIdleConns:        100,              // Maximum idle connections
		MaxIdleConnsPerHost: 10,               // Idle connections per host
		IdleConnTimeout:     90 * time.Second, // Timeout for idle connections
		TLSHandshakeTimeout: 10 * time.Second, // Timeout for TLS handshakes
	},
}
