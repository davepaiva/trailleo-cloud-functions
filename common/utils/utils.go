package utils

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/davepaiva/trailleo-google-cloud-functions/common/types"
	"github.com/google/uuid"
	"googlemaps.github.io/maps"
)
func JsonResponse(w http.ResponseWriter, response types.Response) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    if err := json.NewEncoder(w).Encode(response); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func SetCORSHeaders(w http.ResponseWriter, req *http.Request) bool {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Max-Age", "3600")

	if req.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return true
	}
	return false
}

func GetGoogleMapsToken (sessionTokenStr string) (maps.PlaceAutocompleteSessionToken, error){
	if sessionTokenStr == "" {
		return maps.PlaceAutocompleteSessionToken{}, fmt.Errorf("session_token is required")
	}
	sessionTokenUuid, err := uuid.Parse(sessionTokenStr)
	if err != nil {
		return maps.PlaceAutocompleteSessionToken{}, err
	}
	sessionToken := maps.PlaceAutocompleteSessionToken(sessionTokenUuid)
	return sessionToken, nil
}