package function

import (
	"context"
	"net/http"
	"time"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/davepaiva/trailleo-google-cloud-functions/common/db"
	"github.com/davepaiva/trailleo-google-cloud-functions/common/types"
	"github.com/davepaiva/trailleo-google-cloud-functions/common/utils"
	"go.mongodb.org/mongo-driver/bson"
)

func init(){
	functions.HTTP("GetTrekDetail", GetTrekDetail)
}

func GetTrekDetail (w http.ResponseWriter, req *http.Request){
	slug := req.URL.Query().Get("slug")
	if slug == "" {
		http.Error(w, "Missing slug parameter", http.StatusBadRequest)
		return
	}
	collection := db.Client.Database("trailleo").Collection("treks")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	filter := bson.M{"slug": slug}
	cursor, err := collection.Find(ctx, filter)
	if err!= nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var treks []db.Trek
	cursorErr:= cursor.All(ctx, &treks)
	if cursorErr != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if len(treks) == 0 {
		http.Error(w, "No treks found", http.StatusNotFound)
		return
	}
	response:= types.Response{Message: "success", Data: treks[0]}
	utils.JsonResponse(w, response)
}