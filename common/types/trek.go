package types

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GeoJSON struct {
	Type        string    `bson:"type"`
	Coordinates []float64 `bson:"coordinates"`
}

type AttractionsFilter struct {
	Valley         bool `bson:"Valley"`
	Lakes          bool `bson:"Lakes"`
	PanoramicViews bool `bson:"PanoramicViews"`
	Forest         bool `bson:"Forest"`
	Mountains      bool `bson:"Mountains"`
	WildFlowers    bool `bson:"WildFlowers"`
}

type TrekSummary struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Name       string             `bson:"name" json:"name"`
	PlaceName  string             `bson:"place_name" json:"place_name"`
	CoverImage string             `bson:"cover_image" json:"cover_image"`
	Distance   int                `bson:"distance" json:"distance"`
	Altitude   int                `bson:"altitude" json:"altitude"`
	Difficulty string             `bson:"difficulty" json:"difficulty"`
	Slug       string             `bson:"slug" json:"slug"`
	Duration   string             `bson:"duration" json:"duration"`
}

type TrekSlugItem struct {
	Name       string             `bson:"name" json:"name"`
	Slug       string             `bson:"slug" json:"folder_name"`
}