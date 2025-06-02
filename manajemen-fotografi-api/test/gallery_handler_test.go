package test

import (
	"bytes"
	"encoding/json"
	"log"
	"manajemen-fotografi-api/config"
	"manajemen-fotografi-api/handlers"
	"manajemen-fotografi-api/models"
	"net/http"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestCreateGallery_RealDB(t *testing.T) {
	// Hubungkan ke database MongoDB asli
	config.ConnectDB()

	// Buat Fiber app dan daftarkan handler CreateGallery
	app := fiber.New()
	app.Post("/api/galleries", handlers.CreateGallery)

	// Data dummy gallery
	photographerID := primitive.NewObjectID() // generate ObjectID baru sebagai contoh
	newGallery := models.Gallery{
		PhotographerID: photographerID,
		Title:          "Test Gallery",
		ImageURL:       "https://example.com/image.jpg",
		Description:    "Ini adalah deskripsi galeri test",
		CreatedAt:      time.Now(),
	}
	payload, err := json.Marshal(newGallery)
	if err != nil {
		t.Fatalf("Gagal marshal data gallery: %v", err)
	}

	// Buat request POST dengan JSON payload
	req := bytes.NewReader(payload)
	httpReq, err := http.NewRequest("POST", "/api/galleries", req)
	if err != nil {
		t.Fatalf("Gagal membuat request: %v", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	// Kirim request ke Fiber app
	resp, err := app.Test(httpReq, -1)
	assert.Nil(t, err)
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)

	log.Println("✅ Data gallery berhasil disimpan di MongoDB pada collection `galleries`. Silakan cek MongoDB Compass.")
}
