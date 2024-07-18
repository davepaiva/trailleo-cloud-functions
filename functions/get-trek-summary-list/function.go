package function

import (
	"context"
	"log"
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
	"googlemaps.github.io/maps"
)

func init() {
	functions.HTTP("GetTrekSummaryList", GetTrekSummaryList)
}

func GetTrekSummaryList(w http.ResponseWriter, req *http.Request) {
	if utils.SetCORSHeaders(w, req) {
		return
	}
	search := req.URL.Query().Get("search")
	difficulty_filter := req.URL.Query().Get("difficulty")
	monthsFilter:= req.URL.Query().Get("months")
	whatToExpectFilter:= req.URL.Query().Get("what_to_expect")
	sessionTokenStr:= req.URL.Query().Get("session_token")
	placeId:= req.URL.Query().Get("maps_place_id")

	var (
		lat float64
		lng float64
	)



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

	c, err := maps.NewClient(maps.WithAPIKey(os.Getenv("MAPS_API_KEY")))
	if err!=nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if sessionTokenStr!="" && placeId!=""{
		sessionToken, err := utils.GetGoogleMapsToken(sessionTokenStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		mapsReq := &maps.PlaceDetailsRequest{
			PlaceID:      placeId,
			SessionToken: sessionToken,
		}
	
		mapsRes, err := c.PlaceDetails(context.Background(), mapsReq)
		if err != nil {
			log.Fatalf("fatal error: %s", err)
		}
	
		lat = mapsRes.Geometry.Location.Lat
		lng = mapsRes.Geometry.Location.Lng
	}

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

	    // Add geospatial filter if lat and lng are provided
		if lat != 0 && lng != 0 {
			geoFilter := bson.M{
				"$geoWithin": bson.M{
					"$centerSphere": []interface{}{
						[]float64{lng, lat}, 50.0 / 6378.1, // 50km radius, Earth's radius in km
					},
				},
			}
			pipeline = append(pipeline, bson.M{"$match": bson.M{"location": geoFilter}})
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

	response := types.Response{
		Message: "success",
		Data: trekSummaries,
		Meta: map[string]interface{}{
			"pageCount": totalPages,
			"currentPage": page, // Add current page to the response
		},
	}
	utils.JsonResponse(w, response)
}