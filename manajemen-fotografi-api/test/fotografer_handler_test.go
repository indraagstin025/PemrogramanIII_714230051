package test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http/httptest"
	"time"

	"manajemen-fotografi-api/handlers"
	"manajemen-fotografi-api/models"
	"testing"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)
var photographerCollection *mongo.Collection

func init() {
	// Initialize MongoDB collection (replace with actual DB initialization logic)
	client, err := mongo.NewClient()
	if err != nil {
		panic("Failed to create MongoDB client")
	}
	photographerCollection = client.Database("your_database_name").Collection("photographers")
}
func setupPhotographerApp() *fiber.App {
	app := fiber.New()

	app.Post("/photographers", handlers.CreatePhotographer)
	app.Get("/photographers/:id", handlers.GetPhotographerByID)
	app.Get("/photographers/user/:user_id", handlers.GetPhotographerByUserID)
	app.Get("/photographers", handlers.GetAllPhotographers)
	app.Put("/photographers/:id", handlers.UpdatePhotographer)
	app.Delete("/photographers/:id", handlers.DeletePhotographer)

	return app
}


func TestCreatePhotographerHandler(t *testing.T) {
	app := fiber.New()

	// Pasang route CreatePhotographer
	app.Post("/photographers", handlers.CreatePhotographer)

	// Dummy data input
	userID := primitive.NewObjectID()
	newPhotographer := models.Photographer{
		UserID:      userID,
		Phone:       "08123456789",
		Description: "Fotografer profesional",
		Portfolio:   []string{"http://portfolio.example.com"},
		Location:    "Jakarta",
	}

	// Encode body request JSON
	body, _ := json.Marshal(newPhotographer)

	req := httptest.NewRequest("POST", "/photographers", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusCreated {
		t.Errorf("Expected status %d, got %d", fiber.StatusCreated, resp.StatusCode)
	}

	var created models.Photographer
	if err := json.NewDecoder(resp.Body).Decode(&created); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if created.ID.IsZero() {
		t.Errorf("Expected created photographer to have ID assigned")
	}

	if created.UserID != userID {
		t.Errorf("Expected UserID %s, got %s", userID.Hex(), created.UserID.Hex())
	}
}

func TestGetPhotographerByID(t *testing.T) {
	app := fiber.New()
	app.Get("/photographers/:id", handlers.GetPhotographerByID)

	// Setup: Insert dummy data dulu ke DB supaya bisa dicari
	userID := primitive.NewObjectID()
	dummyPhotographer := models.Photographer{
		ID:          primitive.NewObjectID(),
		UserID:      userID,
		Phone:       "08123456789",
		Description: "Fotografer profesional",
		Portfolio:   []string{"http://portfolio.example.com"},
		Location:    "Jakarta",
		CreatedAt:   time.Now().Unix(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := photographerCollection.InsertOne(ctx, dummyPhotographer)
	if err != nil {
		t.Fatalf("Setup insert failed: %v", err)
	}
	defer func() {
		// Bersihkan data setelah test
		photographerCollection.DeleteOne(ctx, bson.M{"_id": dummyPhotographer.ID})
	}()

	// --- CASE 1: GET dengan ID valid dan ada data ---
	reqValid := httptest.NewRequest("GET", "/photographers/"+dummyPhotographer.ID.Hex(), nil)
	respValid, err := app.Test(reqValid)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	defer respValid.Body.Close()

	if respValid.StatusCode != fiber.StatusOK {
		t.Errorf("Expected status 200, got %d", respValid.StatusCode)
	}

	var result models.Photographer
	if err := json.NewDecoder(respValid.Body).Decode(&result); err != nil {
		t.Fatalf("Decode failed: %v", err)
	}

	if result.ID != dummyPhotographer.ID {
		t.Errorf("Expected ID %s, got %s", dummyPhotographer.ID.Hex(), result.ID.Hex())
	}

	// --- CASE 2: GET dengan ID tidak valid (bukan hex) ---
	reqInvalidID := httptest.NewRequest("GET", "/photographers/invalidhex", nil)
	respInvalidID, err := app.Test(reqInvalidID)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	defer respInvalidID.Body.Close()

	if respInvalidID.StatusCode != fiber.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", respInvalidID.StatusCode)
	}

	// --- CASE 3: GET dengan ID valid tapi tidak ada di DB ---
	nonExistentID := primitive.NewObjectID()
	reqNotFound := httptest.NewRequest("GET", "/photographers/"+nonExistentID.Hex(), nil)
	respNotFound, err := app.Test(reqNotFound)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	defer respNotFound.Body.Close()

	if respNotFound.StatusCode != fiber.StatusNotFound {
		t.Errorf("Expected status 404, got %d", respNotFound.StatusCode)
	}
}






