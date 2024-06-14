package utils

import (
	"encoding/json"
	"net/http"

	"github.com/davepaiva/trailleo-google-cloud-functions/common/types"
)
func JsonResponse(w http.ResponseWriter, response types.Response) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    if err := json.NewEncoder(w).Encode(response); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}