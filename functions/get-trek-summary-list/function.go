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
	"go.mongodb.org/mongo-driver/bson"
)

func init() {
	functions.HTTP("GetTrekSummaryList", GetTrekSummaryList)
}

func GetTrekSummaryList(w http.ResponseWriter, req *http.Request) {
	search := req.URL.Query().Get("search")
	difficulty_filter := req.URL.Query().Get("difficulty")
	monthsFilter:= req.URL.Query().Get("months")
	whatToExpectFilter:= req.URL.Query().Get("what_to_expect")

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
	if difficulty_filter != "" {
		difficultyStage := bson.M{
			"$match": bson.M{
				"difficulty": difficulty_filter,
			},
		}
		pipeline = append(pipeline, difficultyStage)
	}

	if monthsFilter !=""{
		months := strings.Split(monthsFilter, ",")
		monthsFilterStage := bson.M{
			"$match": bson.M{
				"months": bson.M{
					"$in": months,
				},
			},
		}
		pipeline = append(pipeline, monthsFilterStage)
	}

	if whatToExpectFilter != ""{
		whatToExpect:=strings.Split(whatToExpectFilter, ",")
		whatToExpectFilterStage:=bson.M{
			"$match": bson.M{
				"what_to_expect": bson.M{
					"$in": whatToExpect,
				},
			},
		}
		pipeline=append(pipeline, whatToExpectFilterStage)
	}

		// Add a stage to count the total number of documents
		countStage := bson.M{
			"$count": "totalCount",
		}
		countPipeline := append(pipeline, countStage)
	
		countCursor, err := collection.Aggregate(ctx, countPipeline)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	
		var countResult []bson.M
		err = countCursor.All(ctx, &countResult); 
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	
		totalCount := 0
		if len(countResult) > 0 {
			totalCount =  int(countResult[0]["totalCount"].(int32))
		}
	
		totalPages := (totalCount + limit - 1) / limit
	

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
			"slug":       1,
			"duration":   1,
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

	response := types.Response{Message: "success", Data: trekSummaries, Meta: map[string]interface{}{"pageCount": totalPages}};
	utils.JsonResponse(w, response);
}