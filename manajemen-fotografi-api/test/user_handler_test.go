package test

import (
	"bytes"
	"encoding/json"
	"manajemen-fotografi-api/config"
	"manajemen-fotografi-api/handlers"
	"manajemen-fotografi-api/models"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"go.mongodb.org/mongo-driver/bson"
	"context"
	"time"
)

// Setup Fiber app dan middleware session untuk test
func setupApp() *fiber.App {
	app := fiber.New()
	// config.InitMongo() // Pastikan ini menggunakan DB test
	middleware := session.New()

	app.Use(func(c *fiber.Ctx) error {
		sess, _ := middleware.Get(c)
		c.Locals("session", sess)
		return c.Next()
	})

	app.Post("/register", handlers.RegisterUser)
	app.Post("/login", handlers.LoginUser)
	app.Post("/logout", handlers.LogoutUser)

	return app
}

// Hapus user test dari DB agar test bisa diulang tanpa error duplicate
func cleanupTestUser(email string) {
	collection := config.GetCollection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection.DeleteOne(ctx, bson.M{"email": email})
}

func TestRegisterUser(t *testing.T) {
	app := setupApp()

	email := "testuser@example.com"
	cleanupTestUser(email) // hapus dulu kalau ada data lama

	user := map[string]string{
		"name":     "Test User",
		"email":    email,
		"password": "password123",
		"role":     models.RoleClient,
	}

	data, _ := json.Marshal(user)
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusCreated {
		t.Errorf("Expected status 201 Created, got %d", resp.StatusCode)
	}

	// Jangan hapus data supaya bisa dicek di MongoDB Compass
	// cleanupTestUser(email)
}




func TestLogoutUser(t *testing.T) {
	app := setupApp()

	req, _ := http.NewRequest("POST", "/logout", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}
