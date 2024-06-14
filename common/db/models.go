package db

import (
	"time"

	"github.com/davepaiva/trailleo-google-cloud-functions/common/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Trek struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Duration        string             `bson:"duration" json:"duration"`
	Route           string             `bson:"route" json:"route"`
	Polyline        string             `bson:"polyline" json:"polyline"`
	WhatToExpect    []string           `bson:"what_to_expect" json:"what_to_expect"`
	Facilities      []string           `bson:"facilities" json:"facilities"`
	Distance        int                `bson:"distance" json:"distance"`
	CreatorName     string             `bson:"creator_name" json:"creator_name"`
	Video           string             `bson:"video" json:"video"`
	Slug            string             `bson:"slug" json:"slug"`
	CoverImage      string             `bson:"cover_image" json:"cover_image"`
	Altitude        int                `bson:"altitude" json:"altitude"`
	Name            string             `bson:"name" json:"name"`
	CreatorLink     string             `bson:"creator_link" json:"creator_link"`
	Waypoint        string             `bson:"waypoint" json:"waypoint"`
	Location        types.GeoJSON            `bson:"location" json:"location"`
	EntryFee        string             `bson:"entry_fee" json:"entry_fee"`
	BestTime        string             `bson:"best_time" json:"best_time"`
	Description     string             `bson:"description" json:"description"`
	ApproxCost      int                `bson:"approx_cost" json:"approx_cost"`
	ClosestCity     string             `bson:"closest_city" json:"closest_city"`
	FacilityDesc    string             `bson:"facility_desc" json:"facility_desc"`
	GettingThere    string             `bson:"getting_there" json:"getting_there"`
	Difficulty      string             `bson:"difficulty" json:"difficulty"`
	Popularity      int                `bson:"popularity" json:"popularity"`
	UpdatedAt       time.Time          `bson:"updatedAt" json:"updatedAt"`
	Region          string             `bson:"region" json:"region"`
	Months          []string           `bson:"months" json:"months"`
	AdditionalInfo  string             `bson:"additional_info" json:"additional_info"`
	AttractionsFilter types.AttractionsFilter `bson:"attractions_filter" json:"attractions_filter"`
	CreatedAt       time.Time          `bson:"createdAt" json:"createdAt"`
	TrekImages      []string           `bson:"trek_images" json:"trek_images"`
	PlaceName       string             `bson:"place_name" json:"place_name"`
}

