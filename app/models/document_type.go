package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type DocumentType struct {
	ID   primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	Name string             `bson:"name" json:"name"`
}
