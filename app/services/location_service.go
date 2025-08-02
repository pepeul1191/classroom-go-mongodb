package services

import (
	"classroom/app/configs"
	"classroom/app/models"
	"context"
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
