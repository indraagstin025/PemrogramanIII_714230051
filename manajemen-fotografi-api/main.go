package main

import (
	"log"
	"manajemen-fotografi-api/config"
	"manajemen-fotografi-api/middlewares"
	"manajemen-fotografi-api/routes" // Import routes

	"github.com/gofiber/fiber/v2"
)

func main() {
	// Inisialisasi Fiber
	app := fiber.New()

	// Setup middleware

	middlewares.SetupLogger(app)
	middlewares.SetupCORS(app)

	// Koneksi ke database
	config.ConnectDB()

	if config.MongoDatabase == nil {
		log.Fatal("MongoDatabase is nil")
	}

  

	// Setup routes
	routes.SetupRoutes(app) // Menghubungkan semua route yang sudah digabungkan di routes.go

	// Menjalankan server
	log.Fatal(app.Listen(":3000"))
}
