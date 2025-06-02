package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"time"

	"manajemen-fotografi-api/config"
	"manajemen-fotografi-api/models"
	"manajemen-fotografi-api/utils"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var photographerCollection = config.GetCollection("photographers")


// CreatePhotographer menambahkan data fotografer baru
func CreatePhotographer(c *fiber.Ctx) error {
	var photographer models.Photographer

	if err := c.BodyParser(&photographer); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format data salah"})
	}

	if !utils.IsValidPhone(photographer.Phone) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Nomor telepon tidak valid"})
	}


	photographer.CreatedAt = time.Now().Unix()
	photographer.UpdatedAt = photographer.CreatedAt


	// Set ID baru dan timestamp
	photographer.ID = primitive.NewObjectID()
	photographer.CreatedAt = time.Now().Unix()

	// Validasi UserID
	if photographer.UserID.IsZero() {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "User ID harus diisi"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := photographerCollection.InsertOne(ctx, photographer)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal menyimpan fotografer"})
	}

	return c.Status(fiber.StatusCreated).JSON(photographer)
}

// GetPhotographerByID mendapatkan data fotografer berdasarkan ID
func GetPhotographerByID(c *fiber.Ctx) error {
	idParam := c.Params("id")
	photographerID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID tidak valid"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var photographer models.Photographer
	err = photographerCollection.FindOne(ctx, bson.M{"_id": photographerID}).Decode(&photographer)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Fotografer tidak ditemukan"})
	}

	return c.JSON(photographer)
}

// GetPhotographerByUserID mendapatkan data fotografer berdasarkan UserID (relasi)
func GetPhotographerByUserID(c *fiber.Ctx) error {
	userIDParam := c.Params("user_id")
	userID, err := primitive.ObjectIDFromHex(userIDParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "User ID tidak valid"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var photographer models.Photographer
	err = photographerCollection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&photographer)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Fotografer tidak ditemukan"})
	}

	return c.JSON(photographer)
}

// GetAllPhotographers mengambil semua data fotografer
func GetAllPhotographers(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := photographerCollection.Find(ctx, bson.M{})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal mengambil data fotografer"})
	}
	defer cursor.Close(ctx)

	var photographers []models.Photographer
	if err = cursor.All(ctx, &photographers); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal memproses data fotografer"})
	}

	return c.JSON(photographers)
}

// UpdatePhotographer mengubah data fotografer berdasarkan ID
func UpdatePhotographer(c *fiber.Ctx) error {
    idParam := c.Params("id")
    photographerID, err := primitive.ObjectIDFromHex(idParam)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID tidak valid"})
    }

    // ambil file profile_photo jika ada
    file, err := c.FormFile("profile_photo")
    var profilePhotoURL string
    if err == nil {
        // simpan file ke uploads folder
        filename := fmt.Sprintf("%d_%s", time.Now().Unix(), file.Filename)
        err = c.SaveFile(file, fmt.Sprintf("./uploads/%s", filename))
        if err != nil {
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal menyimpan foto profil"})
        }
        profilePhotoURL = fmt.Sprintf("/uploads/%s", filename)
    }

    // ambil data lain dari form fields
    phone := c.FormValue("phone")
    description := c.FormValue("description")
    location := c.FormValue("location")
    portfolioStr := c.FormValue("portfolio") // misal portfolio dikirim sebagai JSON string array dari frontend

    // Validasi nomor telepon jika tidak kosong
    if phone != "" && !utils.IsValidPhone(phone) {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Nomor telepon tidak valid"})
    }

    if len(description) > 500 {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Deskripsi terlalu panjang (maksimal 500 karakter)"})
    }

    if len(location) > 200 {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Lokasi terlalu panjang (maksimal 200 karakter)"})
    }

    // parsing portfolio JSON string ke slice string
    var portfolio []string
    if portfolioStr != "" {
        err = json.Unmarshal([]byte(portfolioStr), &portfolio)
        if err != nil {
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format portfolio salah"})
        }
    }

    // validasi portfolio url tidak kosong
    for _, url := range portfolio {
        if strings.TrimSpace(url) == "" {
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Terdapat URL portfolio yang kosong"})
        }
    }

    update := bson.M{
        "phone":       phone,
        "description": description,
        "portfolio":   portfolio,
        "location":    location,
        "updated_at":  time.Now().Unix(),
    }

    if profilePhotoURL != "" {
        update["profile_photo"] = profilePhotoURL
    }

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    result, err := photographerCollection.UpdateOne(ctx, bson.M{"_id": photographerID}, bson.M{"$set": update})
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal update fotografer"})
    }

    if result.MatchedCount == 0 {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Fotografer tidak ditemukan"})
    }

    var updatedPhotographer models.Photographer
    err = photographerCollection.FindOne(ctx, bson.M{"_id": photographerID}).Decode(&updatedPhotographer)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal mengambil data yang diupdate"})
    }

    return c.JSON(updatedPhotographer)
}



// DeletePhotographer menghapus data fotografer berdasarkan ID
func DeletePhotographer(c *fiber.Ctx) error {
	idParam := c.Params("id")
	photographerID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID tidak valid"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := photographerCollection.DeleteOne(ctx, bson.M{"_id": photographerID})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal menghapus fotografer"})
	}

	if result.DeletedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Fotografer tidak ditemukan"})
	}

	return c.JSON(fiber.Map{"message": "Fotografer berhasil dihapus"})
}