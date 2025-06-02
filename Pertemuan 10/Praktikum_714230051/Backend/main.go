package main

import (
	"fmt"
	"inibackend/config"
	"inibackend/router"
	"log"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("gagal memuat file .env", err)
	}
}

func main() {
	app := fiber.New()

	// Logging Request
	app.Use(logger.New())

	//Basic Cors
	app.Use(cors.New(cors.Config{
		AllowOrigins:     strings.Join(config.GetAllowedOrigins(), ","),
		AllowCredentials: true,
		AllowMethods:     "GET, POST, PUT, DELETE",
	}))
	// Setup application routes
	router.SetupRoutes(app)

	//Handles 404
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  fiber.StatusNotFound,
			"message": "Endpoint Tidak Ditemukan",
		})

	})
	// BACA port yang ada di .env
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000" //default port kalau gak di env
	}

	//untuk log cek konek di port mana
	log.Printf("Server berjalan di port %s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Error Starting server: %v", err)
	} //koneksi terputus
}
