package bdd

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Db *mongo.Client
var Database *mongo.Database

func NewDB() {
	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		uri = "mongodb://localhost:27017" // lancement hors Docker
	}
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "pa2026"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal("connexion MongoDB : ", err)
	}
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal("MongoDB injoignable : ", err)
	}

	Db = client
	Database = client.Database(dbName)
	log.Println("Connexion à MongoDB réussie, base :", dbName)
}