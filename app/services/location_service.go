package services

import (
	"classroom/app/configs"
	"classroom/app/models"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func FetchDepartments() ([]models.LocationMin, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := configs.DB.Collection("locations")

	filter := bson.M{"type": "department"}
	projection := options.Find().SetProjection(bson.M{
		"_id":  1,
		"name": 1,
	})

	cursor, err := collection.Find(ctx, filter, projection)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []models.LocationMin
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}

func FetchProvincesByDepartment(departmentID string) ([]models.LocationMin, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := configs.DB.Collection("locations")

	// Convertir el string ID a ObjectID
	objID, err := primitive.ObjectIDFromHex(departmentID)
	if err != nil {
		return nil, err
	}

	// Filtro: provincias con parent_id = departamento
	filter := bson.M{
		"type":      "province",
		"parent_id": objID,
	}

	// Proyección mínima
	opts := options.Find().SetProjection(bson.M{
		"_id":  1,
		"name": 1,
	})

	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []models.LocationMin
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}

func FetchDistrictsByProvince(provinceID string) ([]models.LocationMin, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := configs.DB.Collection("locations")

	// Convertir el string ID a ObjectID
	objID, err := primitive.ObjectIDFromHex(provinceID)
	if err != nil {
		return nil, err
	}

	// Filtro: provincias con parent_id = departamento
	filter := bson.M{
		"type":      "district",
		"parent_id": objID,
	}

	// Proyección mínima
	opts := options.Find().SetProjection(bson.M{
		"_id":  1,
		"name": 1,
	})

	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []models.LocationMin
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}
