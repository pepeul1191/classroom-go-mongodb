package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func mai2n() {
	// Crear un cliente de MongoDB
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017") // Cambia la URI según tu configuración
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.Background())

	// Verificar la conexión
	if err := client.Ping(context.Background(), readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Conexión exitosa a MongoDB!")
	// name := "La V"

	// Seleccionar la base de datos y colección
	db := client.Database("peru")            // Reemplaza "mydb" con el nombre de tu base de datos
	collection := db.Collection("locations") // Reemplaza "locations" con el nombre de tu colección

	// Construir la consulta
	pipeline := []bson.M{
		{
			"$match": bson.M{
				"type": "district",
			},
		},
		{
			"$graphLookup": bson.M{
				"from":             "locations",
				"startWith":        "$parent_id",
				"connectFromField": "parent_id",
				"connectToField":   "_id",
				"as":               "ancestors",
			},
		},
		{
			"$addFields": bson.M{
				"full_name": bson.M{
					"$concat": []interface{}{
						"$name",
						", ",
						bson.M{
							"$first": bson.M{
								"$map": bson.M{
									"input": bson.M{
										"$filter": bson.M{
											"input": "$ancestors",
											"as":    "a",
											"cond":  bson.M{"$eq": []interface{}{"$$a.type", "province"}},
										},
									},
									"as": "prov",
									"in": "$$prov.name",
								},
							},
						},
						", ",
						bson.M{
							"$first": bson.M{
								"$map": bson.M{
									"input": bson.M{
										"$filter": bson.M{
											"input": "$ancestors",
											"as":    "a",
											"cond":  bson.M{"$eq": []interface{}{"$$a.type", "department"}},
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
		{
			"$match": bson.M{
				"full_name": bson.M{
					"$regex":   "La V", // Pasando directamente el string para el regex
					"$options": "i",    // Opciones para la búsqueda insensible a mayúsculas/minúsculas
				},
			},
		},
		{
			"$project": bson.M{
				"_id":         0,
				"district_id": bson.M{"$toString": "$_id"},
				"full_name":   1,
			},
		},
	}

	// Ejecutar la consulta
	cursor, err := collection.Aggregate(context.Background(), pipeline)
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.Background())

	// Iterar sobre los resultados
	for cursor.Next(context.Background()) {
		var result bson.M
		if err := cursor.Decode(&result); err != nil {
			log.Fatal(err)
		}
		fmt.Println(result)
	}

	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}
}
