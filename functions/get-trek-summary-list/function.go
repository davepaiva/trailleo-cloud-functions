package function

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/davepaiva/trailleo-google-cloud-functions/common/db"
	"github.com/davepaiva/trailleo-google-cloud-functions/common/types"
	"github.com/davepaiva/trailleo-google-cloud-functions/common/utils"
	"go.mongodb.org/mongo-driver/bson"
)

func init() {
	functions.HTTP("GetTrekSummaryList", GetTrekSummaryList)
}

func GetTrekSummaryList(w http.ResponseWriter, req *http.Request) {
	search := req.URL.Query().Get("search")
	difficulty := req.URL.Query().Get("difficulty")

	page, err := strconv.Atoi(req.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(req.URL.Query().Get("limit"))
	if err != nil || limit < 1 {
		limit = 10
	}
	skip := (page - 1) * limit
	collection := db.Client.Database("trailleo").Collection("treks")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Initialize the pipeline for running multiple aggregations= filters
	pipeline := []bson.M{}

	// Create the search filter using Atlas Search autocomplete with compound should and scoring
	if search != "" {
		searchStage := bson.M{
			"$search": bson.M{
				"index": "trek_text_location_autocomplete", // Name of your Atlas Search index
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

	// Create the difficulty filter
	if difficulty != "" {
		difficultyStage := bson.M{
			"$match": bson.M{
				"difficulty": difficulty,
			},
		}
		pipeline = append(pipeline, difficultyStage)
	}

	// Add skip, limit, and project stages
	pipeline = append(pipeline,
		bson.M{"$skip": skip},
		bson.M{"$limit": limit},
		bson.M{"$project": bson.M{
			"_id":        1,
			"name":       1,
			"place_name": 1,
			"cover_image": 1,
			"distance":   1,
			"altitude":   1,
			"difficulty": 1,
		}},
	)

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var trekSummaries []types.TrekSummary
	cursorErr := cursor.All(ctx, &trekSummaries)
	if cursorErr != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := types.Response{Message: "success", Data: trekSummaries}
	utils.JsonResponse(w, response)
}