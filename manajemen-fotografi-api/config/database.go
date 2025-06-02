package config

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var MongoClient *mongo.Client
var MongoDatabase *mongo.Database

// ConnectDB untuk menghubungkan ke database MongoDB
func ConnectDB() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: Error loading .env file:", err)
		// Continue execution, as environment variables might be set in the system
	}

	// Get MongoDB URI from environment variable
	mongoURI := os.Getenv("MONGOSTRING")
	if mongoURI == "" {
		log.Fatal("MONGOSTRING is not set in environment variables")
	}
	log.Println("Mongo URI loaded successfully")

	// Set context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create MongoDB client
	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}

	// Ping MongoDB to verify connection
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal("Failed to ping MongoDB:", err)
	}
	log.Println("Successfully connected to MongoDB")

	// Initialize MongoClient and MongoDatabase
	MongoClient = client
	MongoDatabase = client.Database("manajemen-fotografi")
	log.Println("MongoDB database 'manajemen-fotografi' is initialized successfully")
}

// GetCollection untuk mengambil koleksi tertentu berdasarkan nama
func GetCollection(collectionName string) *mongo.Collection {
	if MongoDatabase == nil {
		// Ensure database is connected if GetCollection is called before ConnectDB
		log.Println("Database connection not initialized, connecting now...")
		ConnectDB()
		
		if MongoDatabase == nil {
			log.Fatal("Failed to initialize database connection")
		}
	}
	return MongoDatabase.Collection(collectionName)
}

// DisconnectDB untuk memutuskan koneksi dengan database MongoDB
func DisconnectDB() {
	if MongoClient != nil {
		err := MongoClient.Disconnect(context.Background())
		if err != nil {
			log.Println("Error disconnecting from MongoDB:", err)
		} else {
			log.Println("Successfully disconnected from MongoDB")
		}
	}
}