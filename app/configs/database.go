package configs

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database

func ConnectToMongoDB() error {
	err := godotenv.Load()
	if err != nil {
		log.Println("No se pudo cargar .env, usando variables del sistema")
	}

	mongoURI := os.Getenv("MONGODB_URI")
	dbName := os.Getenv("MONGODB_DATABASE")

	if mongoURI == "" || dbName == "" {
		log.Fatal("Faltan variables de entorno: MONGODB_URI o MONGODB_DATABASE")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOpts := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return err
	}

	// Verificar conexión
	if err := client.Ping(ctx, nil); err != nil {
		return err
	}

	DB = client.Database(dbName)
	log.Println("✅ Conectado a MongoDB")
	return nil
}
