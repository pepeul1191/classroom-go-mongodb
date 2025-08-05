package services

import (
	"classroom/app/configs"
	"classroom/app/models"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// interfáz

type DocumentTypesService interface {
	FetchAll() ([]models.DocumentType, error)
	FetchOne(ID *primitive.ObjectID) (*models.DocumentType, error)
}

type documentTypesServiceImpl struct{}

func NewDocumentTypesService() DocumentTypesService {
	return &documentTypesServiceImpl{}
}

// métodos

func (s *documentTypesServiceImpl) FetchAll() ([]models.DocumentType, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := configs.DB.Collection("document_types")

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []models.DocumentType
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}

func (s *documentTypesServiceImpl) FetchOne(ID *primitive.ObjectID) (*models.DocumentType, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := configs.DB.Collection("document_types")

	var result models.DocumentType
	err := collection.FindOne(ctx, bson.M{"_id": ID}).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
