package function

import (
	"context"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/davepaiva/trailleo-google-cloud-functions/common/db"
	"github.com/davepaiva/trailleo-google-cloud-functions/common/types"
	"github.com/davepaiva/trailleo-google-cloud-functions/common/utils"
)

func init() {
	functions.HTTP("GetTrekFileList", GetTrekFileList)
}

func GetTrekFileList (w http.ResponseWriter, req *http.Request){
	collection:= db.Client.Database("trailleo").Collection("treks")
	ctx, cancel:= context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	projection := bson.M{
	  "name":       1,
	  "slug": 	    1,
  }
	cursor, err:= collection.Find(ctx, bson.M{}, options.Find().SetProjection(projection))
	if err!=nil {
	  http.Error(w, "Internal error", http.StatusInternalServerError)
	}
	var trekSlugs []types.TrekSlugItem
	cursorErr:= cursor.All(ctx, &trekSlugs)
	if cursorErr!=nil{
	  http.Error(w, "Internal error", http.StatusInternalServerError)
	}
	response:= types.Response{Message: "success", Data: trekSlugs}
	utils.JsonResponse(w, response)
}

