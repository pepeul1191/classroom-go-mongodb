package services

import (
	"classroom/app/configs"
	"classroom/app/models"
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// interfáz

type LocationsService interface {
	FetchDepartments() ([]models.LocationMin, error)
	FetchProvincesByDepartment(departmentID string) ([]models.LocationMin, error)
	FetchDistrictsByProvince(provinceID string) ([]models.LocationMin, error)
	FindDistrictsByFullName(name string, limit uint) ([]models.LocationResult, error)
	InsertDepartment(dep models.LocationMin) (*models.Location, error)
	InsertProvince(pro models.LocationMin, deparmentId primitive.ObjectID) (*models.Location, error)
	ProcessLocations(news []models.NewLocation, edits []models.EditLocation, deletes []primitive.ObjectID, loctionType string, parentId *primitive.ObjectID) ([]models.CreatedLocationResponse, error)
}

type locationsServiceImpl struct{}

func NewLocationsService() LocationsService {
	return &locationsServiceImpl{}
}

// métodos

func (s *locationsServiceImpl) FetchDepartments() ([]models.LocationMin, error) {
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

func (s *locationsServiceImpl) FetchProvincesByDepartment(departmentID string) ([]models.LocationMin, error) {
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

func (s *locationsServiceImpl) FetchDistrictsByProvince(provinceID string) ([]models.LocationMin, error) {
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

func (s *locationsServiceImpl) FindDistrictsByFullName(name string, limit uint) ([]models.LocationResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := configs.DB.Collection("locations")

	pipeline := mongo.Pipeline{
		// Solo distritos
		{{Key: "$match", Value: bson.M{"type": "district"}}},

		// Buscar ancestros
		{{
			Key: "$graphLookup",
			Value: bson.M{
				"from":             "locations",
				"startWith":        "$parent_id",
				"connectFromField": "parent_id",
				"connectToField":   "_id",
				"as":               "ancestors",
			},
		}},

		// Construir nombre completo
		{{
			Key: "$addFields",
			Value: bson.M{
				"full_name": bson.M{
					"$concat": []interface{}{
						"$name", ", ",
						bson.M{
							"$first": bson.A{
								bson.M{
									"$map": bson.M{
										"input": bson.M{
											"$filter": bson.M{
												"input": "$ancestors",
												"as":    "a",
												"cond":  bson.M{"$eq": bson.A{"$$a.type", "province"}},
											},
										},
										"as": "prov",
										"in": "$$prov.name",
									},
								},
							},
						}, ", ",
						bson.M{
							"$first": bson.A{
								bson.M{
									"$map": bson.M{
										"input": bson.M{
											"$filter": bson.M{
												"input": "$ancestors",
												"as":    "a",
												"cond":  bson.M{"$eq": bson.A{"$$a.type", "department"}},
											},
										},
										"as": "dep",
										"in": "$$dep.name",
									},
								},
							},
						},
					},
				},
			},
		}},

		// Filtro tipo LIKE
		{{
			Key: "$match",
			Value: bson.M{
				"full_name": bson.M{"$regex": name, "$options": "i"},
			},
		}},

		// Proyección final
		{{
			Key: "$project",
			Value: bson.M{
				"_id":         0,
				"district_id": bson.M{"$toString": "$_id"},
				"full_name":   1,
			},
		}},
		{{Key: "$limit", Value: limit}},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []models.LocationResult
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}

func (s *locationsServiceImpl) InsertDepartment(dep models.LocationMin) (*models.Location, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := configs.DB.Collection("locations")

	now := time.Now()

	newDepartment := models.Location{
		ID:      primitive.NewObjectID(),
		Name:    dep.Name,
		Type:    "department",
		Created: now,
		Updated: now,
	}

	_, err := collection.InsertOne(ctx, newDepartment)
	if err != nil {
		return nil, err
	}

	return &newDepartment, nil
}

func (s *locationsServiceImpl) InsertProvince(loc models.LocationMin, deparmentId primitive.ObjectID) (*models.Location, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := configs.DB.Collection("locations")

	now := time.Now()

	newProvince := models.Location{
		ID:       loc.ID,
		Name:     loc.Name,
		Type:     "province",
		Created:  now,
		Updated:  now,
		ParentID: &deparmentId,
	}

	_, err := collection.InsertOne(ctx, newProvince)
	if err != nil {
		return nil, err
	}

	return &newProvince, nil
}

func (s *locationsServiceImpl) ProcessLocations(news []models.NewLocation, edits []models.EditLocation, deletes []primitive.ObjectID, loctionType string, parentId *primitive.ObjectID) ([]models.CreatedLocationResponse, error) {
	response := make([]models.CreatedLocationResponse, 0)
	ctx := context.Background()
	collection := configs.DB.Collection("locations")

	// 1. Crear nuevas ubicaciones (sin transacción)
	for _, incoming := range news {
		newLocation := models.Location{
			Name:     incoming.Name,
			Type:     loctionType,
			Created:  time.Now(),
			Updated:  time.Now(),
			ParentID: parentId,
		}

		result, err := collection.InsertOne(ctx, newLocation)
		if err != nil {
			return nil, fmt.Errorf("error creando ubicación: %v", err)
		}

		response = append(response, models.CreatedLocationResponse{
			Tmp: incoming.ID,
			ID:  result.InsertedID.(primitive.ObjectID).Hex(),
		})
		debug := func(msg string, data interface{}) {
			fmt.Printf("[DEBUG] %s: %+v\n", msg, data)
			log.Printf("[LOG] %s: %+v", msg, data)
		}
		debug("Después del append", response)
	}

	// 2. Actualizar (sin transacción)
	for _, incoming := range edits {
		objID, err := primitive.ObjectIDFromHex(incoming.ID)
		if err != nil {
			return nil, fmt.Errorf("ID inválido: %s", incoming.ID)
		}

		_, err = collection.UpdateOne(
			ctx,
			bson.M{"_id": objID},
			bson.M{"$set": bson.M{
				"name":    incoming.Name,
				"updated": time.Now(),
			}},
		)
		if err != nil {
			return nil, fmt.Errorf("error actualizando ubicación %s: %v", incoming.ID, err)
		}
	}

	// 3. Eliminar (sin transacción)
	if len(deletes) > 0 {
		_, err := collection.DeleteMany(
			ctx,
			bson.M{"_id": bson.M{"$in": deletes}},
		)
		if err != nil {
			return nil, fmt.Errorf("error eliminando ubicaciones: %v", err)
		}
	}

	return response, nil
}
