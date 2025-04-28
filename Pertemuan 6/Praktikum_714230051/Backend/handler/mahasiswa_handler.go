package handler

import (
	"inibackend/repository"
	"inibackend/model"
	"github.com/gofiber/fiber/v2"

	"strconv"
)

// Homepage handler
func Homepage(c *fiber.Ctx) error {
	return c.SendString("Welcome to the homepage!")
}

// GetAllMahasiswa handler
func GetAllMahasiswa(c *fiber.Ctx) error {
	data, err := repository.GetAllMahasiswa(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal Mengambil Data Mahasiswa",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  fiber.StatusOK,
		"message": "Berhasil Mengambil Data Mahasiswa",
		"data":    data,
	})
}



// GetAllMahasiswaByNPM handler
func GetAllMahasiswaByNPM(c *fiber.Ctx) error {
	npmStr := c.Params("npm")
	npm, err := strconv.Atoi(npmStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "NPM harus berupa angka",
		})
	}

	mhs, err := repository.GetMahasiswaByNPM(c.Context(), npm)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if mhs == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Mahasiswa tidak ditemukan",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  fiber.StatusOK,
		"message": "Berhasil Mengambil Data Mahasiswa",
		"data":    mhs,
	})
}

// CreateMahasiswa handler (POST)
func CreateMahasiswa(c *fiber.Ctx) error {
	var mhs model.Mahasiswa
	if err := c.BodyParser(&mhs); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Gagal Parse Body",
		})
	}

	// Insert Mahasiswa
	insertedID, err := repository.InsertMahasiswa(c.Context(), mhs)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":      fiber.StatusOK,
		"message":     "Berhasil Menambahkan Mahasiswa",
		"inserted_id": insertedID,
	})
}

// UpdateMahasiswa handler (PUT)
func UpdateMahasiswa(c *fiber.Ctx) error {
	npmStr := c.Params("npm")
	npm, err := strconv.Atoi(npmStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "NPM harus berupa angka",
		})
	}

	var mhs model.Mahasiswa
	if err := c.BodyParser(&mhs); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Gagal Parse Body",
		})
	}

	// Update Mahasiswa
	updatedNPM, err := repository.UpdateMahasiswa(c.Context(), npm, mhs)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  fiber.StatusOK,
		"message": "Berhasil Mengupdate Mahasiswa",
		"npm":     updatedNPM,
	})
}

// DeleteMahasiswa handler (DELETE)
func DeleteMahasiswa(c *fiber.Ctx) error {
	npmStr := c.Params("npm")
	npm, err := strconv.Atoi(npmStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "NPM harus berupa angka",
		})
	}

	// Delete Mahasiswa
	deletedNPM, err := repository.DeleteMahasiswa(c.Context(), npm)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  fiber.StatusOK,
		"message": "Berhasil Menghapus Mahasiswa",
		"npm":     deletedNPM,
	})
}
