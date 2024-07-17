package function

import (
	"context"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/davepaiva/trailleo-google-cloud-functions/common/db"
	"github.com/davepaiva/trailleo-google-cloud-functions/common/types"
	"github.com/davepaiva/trailleo-google-cloud-functions/common/utils"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"googlemaps.github.io/maps"
)

func init() {
	functions.HTTP("GetSearchSuggestions", GetSearchSuggestions)
}

type CombinedResponse struct {
    AutocompleteResults []types.GmapAutoCompleteResponseItem   `json:"autocomplete_results"`
    DbTrekSearchResults       []types.DbTrekSearchResultsItem  `json:"db_trek_search_results"`
}

func GetSearchSuggestions(w http.ResponseWriter, req *http.Request) {
	if utils.SetCORSHeaders(w, req) {
		return
	}
	search := req.URL.Query().Get("search")

	sessionTokenStr := req.URL.Query().Get("session_token")
	if sessionTokenStr == "" {
		http.Error(w, "session_token is required", http.StatusBadRequest)
		return
	}

	sessionTokenUuid, err := uuid.Parse(sessionTokenStr)
	if err != nil {
		http.Error(w, "Invalid session_token format", http.StatusBadRequest)
		return
	}

	sessionToken := maps.PlaceAutocompleteSessionToken(sessionTokenUuid)

	page, err := strconv.Atoi(req.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(req.URL.Query().Get("limit"))
	if err != nil || limit < 1 {
		limit = 10
	}
	skip := (page - 1) * limit
	collection := db.Client.Database(os.Getenv("DB_NAME")).Collection("treks")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Initialize the pipeline for running multiple aggregations= filters
	pipeline := []bson.M{}

	// Create the search filter using Atlas Search autocomplete with compound should and scoring
	if search != "" {
		searchStage := bson.M{
			"$search": bson.M{
				"index": "trek_name_place_autocomplete", // Name of your Atlas Search index
				"compound": bson.M{
					"should": []bson.M{
						{
							"autocomplete": bson.M{
								"query": search,
								"path":  "name",
								"score": bson.M{"boost": bson.M{"value": 10}},
							},
						},
						{
							"autocomplete": bson.M{
								"query": search,
								"path":  "place_name",
								"score": bson.M{"boost": bson.M{"value": 5}},
							},
						},
						{
							"autocomplete": bson.M{
								"query": search,
								"path":  "closest_city",
								"score": bson.M{"boost": bson.M{"value": 1}},
							},
						},
					},
				},
			},
		}
		pipeline = append(pipeline, searchStage)
	}
	

	// Add skip, limit, and project stages
	pipeline = append(pipeline,
		bson.M{"$skip": skip},
		bson.M{"$limit": limit},
		bson.M{"$project": bson.M{
			"_id":        1,
			"name":       1,
			"place_name": 1,
		}},
	)

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var dbTrekSearchResults []types.DbTrekSearchResultsItem
	cursorErr := cursor.All(ctx, &dbTrekSearchResults)
	if cursorErr != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	r := &maps.PlaceAutocompleteRequest{
		Input: search,
		Components: map[maps.Component][]string{
			maps.ComponentCountry: {"in"},
		},
		Types: maps.AutocompletePlaceTypeGeocode,
		SessionToken: sessionToken,
	}
	c, err := maps.NewClient(maps.WithAPIKey(os.Getenv("MAPS_API_KEY")))
	if err!=nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	googleAutocompleteResult, err:=c.PlaceAutocomplete(context.Background(), r)
	if err!=nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var autocompleteResults []types.GmapAutoCompleteResponseItem
	for _, prediction := range googleAutocompleteResult.Predictions {
		description:= prediction.Description
		lastCommaIndex:= strings.LastIndex(description, ",")
		if lastCommaIndex!= -1 {
			description = strings.TrimSpace(description[:lastCommaIndex])
		}
		autocompleteResults = append(autocompleteResults, types.GmapAutoCompleteResponseItem{
			Type: "places",
			Id: prediction.PlaceID,
			Name: prediction.StructuredFormatting.MainText,
			PlaceName: description,
		})
	}

	combinedData := CombinedResponse{
		DbTrekSearchResults: dbTrekSearchResults,
		AutocompleteResults: autocompleteResults,
	}


	response := types.Response{
		Message: "success",
		Data: combinedData,
		Meta: map[string] string{
			"session_token": sessionTokenStr,
		},
	}
	utils.JsonResponse(w, response)
}