package google

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Nick-Anderssohn/oidc-demo/internal/oidc"
)

const googleDiscoveryURL = "https://accounts.google.com/.well-known/openid-configuration"

func GetDiscoveryData(w http.ResponseWriter, r *http.Request) {
	discoveryData, err := oidc.GetDiscoveryData(googleDiscoveryURL)
	if err != nil {
		log.Println("Error fetching discovery data:", err)
		http.Error(w, "Failed to fetch discovery data", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(discoveryData); err != nil {
		log.Println("Error encoding discovery data:", err)
		http.Error(w, "Failed to encode discovery data", http.StatusInternalServerError)
		return
	}
}

func getGoogleDiscoveryData() (*oidc.DiscoveryData, error) {
	return oidc.GetDiscoveryData(googleDiscoveryURL)
}
