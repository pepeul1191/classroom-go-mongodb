package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Person base structure (embedded in other types)
type Person struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	Type           string             `bson:"type" json:"type"` // "teacher", "student", "representative"
	Names          string             `bson:"names" json:"names"`
	LastNames      string             `bson:"last_names" json:"last_names"`
	DocumentNumber string             `bson:"document_number" json:"document_number"`
	ImageURL       string             `bson:"image_url" json:"image_url"`
	Addresses      []Address          `bson:"addresses" json:"addresses"`
	Phones         []Phone            `bson:"phones" json:"phones"`
	User           User               `bson:"user" json:"user"`
	Created        time.Time          `bson:"created" json:"created"`
	Updated        time.Time          `bson:"updated" json:"updated"`
	DocumentTypeID primitive.ObjectID `bson:"document_type_id" json:"document_type_id"`
}

type Address struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	Name       string             `bson:"name"`
	LocationID primitive.ObjectID `bson:"location_id"`
	Created    time.Time          `bson:"created"`
	Updated    time.Time          `bson:"updated"`
}

type Phone struct {
	ID      primitive.ObjectID `bson:"_id,omitempty"`
	Name    string             `bson:"name"`
	Number  string             `bson:"name"`
	Created time.Time          `bson:"created"`
	Updated time.Time          `bson:"updated"`
}

// Teacher specific fields
type Teacher struct {
	Person `bson:",inline"`
	Code   string `bson:"code" json:"code"`
}

type TeacherCreateRequest struct {
	Names          string `json:"names" validate:"required"`
	LastNames      string `json:"last_names" validate:"required"`
	ImageURL       string `json:"image_url" validate:"required"`
	DocumentNumber string `json:"document_number" validate:"required"`
	DocumentTypeID string `json:"document_type_id" validate:"required"`
	Email          string `json:"email" validate:"required,email"`
	Code           string `json:"code" validate:"required"`
}

type TeacherResponse struct {
	ID             string    `json:"_id"`
	Names          string    `json:"names"`
	LastNames      string    `json:"last_names"`
	DocumentNumber string    `json:"document_number"`
	Code           string    `json:"code"`
	Email          string    `json:"email"`
	CreatedAt      time.Time `json:"created"`
}

// Student specific fields
type Student struct {
	Person          `bson:",inline"`
	Code            string               `bson:"code"`
	Representatives []RepresentativeRole `bson:"representatives"`
}

// Representative specific fields
type Representative struct {
	Person `bson:",inline"`
	// Puedes añadir campos específicos aquí
}

// RepresentativeRole embedded
type RepresentativeRole struct {
	ID               primitive.ObjectID `bson:"_id"`
	Relation         string             `bson:"relation"`
	RepresentativeID primitive.ObjectID `bson:"representative_id"`
}
