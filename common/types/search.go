package types

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GmapAutoCompleteResponseItem struct {
	Type string     `json:"type"`
	Id string       `json:"id"`
	Name string     `json:"name"`
	PlaceName string `json:"place_name"`
	
}

type DbTrekSearchResultsItem struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Name       string             `bson:"name" json:"name"`
	PlaceName  string             `bson:"place_name" json:"place_name"`
}