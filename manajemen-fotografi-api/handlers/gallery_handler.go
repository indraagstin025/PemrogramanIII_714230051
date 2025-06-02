package handlers

import (
	"context"
	"time"

	"manajemen-fotografi-api/config"
	"manajemen-fotografi-api/models"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var galleryCollection = config.GetCollection("galleries")

func GetAllGalleries(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := galleryCollection.Find(ctx, bson.M{})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal mengambil galeri"})
	}
	defer cursor.Close(ctx)

	var galleries []models.Gallery
	if err := cursor.All(ctx, &galleries); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal decode data galeri"})
	}

	return c.JSON(galleries)
}

func GetGalleryByID(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID tidak valid"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var gallery models.Gallery
	err = galleryCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&gallery)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Galeri tidak ditemukan"})
	}

	return c.JSON(gallery)
}

func CreateGallery(c *fiber.Ctx) error {
	var gallery models.Gallery

	if err := c.BodyParser(&gallery); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format data salah"})
	}

	if gallery.PhotographerID.IsZero() || gallery.Title == "" || gallery.ImageURL == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Photographer ID, Title, dan ImageURL wajib diisi"})
	}

	gallery.ID = primitive.NewObjectID()
	gallery.CreatedAt = time.Now()
	gallery.UpdatedAt = gallery.CreatedAt

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := galleryCollection.InsertOne(ctx, gallery)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal menyimpan galeri"})
	}

	return c.Status(fiber.StatusCreated).JSON(gallery)
}


func UpdateGallery(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID tidak valid"})
	}

	var updated models.Gallery
	if err := c.BodyParser(&updated); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format data salah"})
	}

	updated.UpdatedAt = time.Now()

	update := bson.M{
		"$set": bson.M{
			"photographer_id": updated.PhotographerID,
			"title":           updated.Title,
			"image_url":       updated.ImageURL,
			"description":     updated.Description,
			"updated_at":      updated.UpdatedAt,
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := galleryCollection.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal update galeri"})
	}
	if result.MatchedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Galeri tidak ditemukan"})
	}

	return c.JSON(fiber.Map{"message": "Galeri berhasil diperbarui"})
}

func DeleteGallery(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID tidak valid"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := galleryCollection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal hapus galeri"})
	}
	if result.DeletedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Galeri tidak ditemukan"})
	}

	return c.JSON(fiber.Map{"message": "Galeri berhasil dihapus"})
}
