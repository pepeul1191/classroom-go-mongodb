package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Location struct {
	ID       primitive.ObjectID  `bson:"_id,omitempty" json:"id"`
	Name     string              `bson:"name" json:"name"`
	Type     string              `bson:"type" json:"type"`
	ParentID *primitive.ObjectID `bson:"parent_id,omitempty" json:"parent_id,omitempty"`
	Created  time.Time           `bson:"created,omitempty" json:"created"`
	Updated  time.Time           `bson:"updated,omitempty" json:"updated"`
}

type LocationMin struct {
	ID   primitive.ObjectID `bson:"_id" json:"_id"`
	Name string             `bson:"name" json:"name"`
}

type LocationResult struct {
	DistrictID string `bson:"district_id" json:"district_id"`
	FullName   string `bson:"full_name" json:"full_name"`
}
