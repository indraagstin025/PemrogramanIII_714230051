package test

import (
	"bytes"
	"context"
	"encoding/json"
	"manajemen-fotografi-api/config"
	"manajemen-fotografi-api/handlers"
	"manajemen-fotografi-api/models"
	"net/http"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Setup Fiber app dan middleware session untuk test client
func setupClientApp() *fiber.App {
	app := fiber.New()
	// config.InitMongo() // Removed because InitMongo is undefined
	middleware := session.New()

	app.Use(func(c *fiber.Ctx) error {
		sess, _ := middleware.Get(c)
		c.Locals("session", sess)
		return c.Next()
	})

	app.Post("/clients", handlers.CreateClient)
	app.Get("/clients/:id", handlers.GetClientByID)
	app.Get("/clients/user/:user_id", handlers.GetClientByUserID)
	app.Put("/clients/:id", handlers.UpdateClient)

	return app
}

func TestCreateClient(t *testing.T) {
	app := setupClientApp()

	// Data client dummy dengan UserID valid dummy
	userID := primitive.NewObjectID()

	client := models.Client{
		UserID: userID,
		Phone:  "081234567890",
	}
	data, _ := json.Marshal(client)

	req := bytes.NewBuffer(data)
	httpReq, err := http.NewRequest("POST", "/clients", req)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(httpReq)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	if resp.StatusCode != fiber.StatusCreated {
		t.Errorf("Expected status 201, got %d", resp.StatusCode)
	}
}

func TestGetClientByID(t *testing.T) {
	app := setupClientApp()

	// Buat dulu client untuk diambil datanya
	client := models.Client{
		ID:     primitive.NewObjectID(),
		UserID: primitive.NewObjectID(),
		Phone:  "081234567890",
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := config.GetCollection("clients").InsertOne(ctx, client)
	if err != nil {
		t.Fatalf("Failed insert client: %v", err)
	}

	url := "/clients/" + client.ID.Hex()
	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	resp, err := app.Test(httpReq)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestGetClientByUserID(t *testing.T) {
	app := setupClientApp()

	// Buat dulu client untuk diambil data berdasarkan user_id
	userID := primitive.NewObjectID()
	client := models.Client{
		ID:     primitive.NewObjectID(),
		UserID: userID,
		Phone:  "081234567890",
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	_, err := config.GetCollection("clients").InsertOne(ctx, client)
	if err != nil {
		t.Fatalf("Failed insert client: %v", err)
	}

	url := "/clients/user/" + userID.Hex()
	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	
	resp, err := app.Test(httpReq)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}


func TestUpdateClient(t *testing.T) {
	app := setupClientApp()

	// Buat dulu client yang akan diupdate
	client := models.Client{
		ID:     primitive.NewObjectID(),
		UserID: primitive.NewObjectID(),
		Phone:  "081234567890",
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := config.GetCollection("clients").InsertOne(ctx, client)
	if err != nil {
		t.Fatalf("Failed insert client: %v", err)
	}

	// Update data
	updatedClient := models.Client{
		Phone: "089876543210",
	}
	data, _ := json.Marshal(updatedClient)
	url := "/clients/" + client.ID.Hex()
	req := bytes.NewBuffer(data)
	httpReq, err := http.NewRequest("PUT", url, req)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(httpReq)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}
