package router

import (
	"inibackend/handler"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")

	api.Get("/", handler.Homepage)
	api.Get("/mahasiswa", handler.GetAllMahasiswa)
	api.Get("/mahasiswa/:npm", handler.GetAllMahasiswaByNPM)
	api.Post("/mahasiswa", handler.CreateMahasiswa)     // POST untuk tambah mahasiswa
	api.Put("/mahasiswa/:npm", handler.UpdateMahasiswa)  // PUT untuk update mahasiswa
	api.Delete("/mahasiswa/:npm", handler.DeleteMahasiswa)
	

	 // DELETE untuk hapus mahasiswa
}
