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

var bookingHandlerCollection = config.GetBookingCollection()

func isValidBookingStatus(status string) bool {
    return status == models.BookingStatusPending ||
           status == models.BookingStatusConfirmed ||
           status == models.BookingStatusDone
}


func GetAllBookings(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := bookingHandlerCollection.Find(ctx, bson.M{})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal mengambil data"})
	}
	defer cursor.Close(ctx)

	var bookings []models.Booking
	if err := cursor.All(ctx, &bookings); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal decode data"})
	}

	return c.JSON(bookings)
}

func GetBookingByID(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID tidak valid"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var booking models.Booking
	err = bookingHandlerCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&booking)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Booking tidak ditemukan"})
	}

	return c.JSON(booking)
}

func CreateBooking(c *fiber.Ctx) error {
	var booking models.Booking

	if err := c.BodyParser(&booking); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format data salah"})
	}

	// Validasi status booking
	if !isValidBookingStatus(booking.Status) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Status booking tidak valid"})
	}

	booking.ID = primitive.NewObjectID()
	booking.CreatedAt = time.Now()
	booking.UpdatedAt = booking.CreatedAt

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := bookingHandlerCollection.InsertOne(ctx, booking)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal menyimpan booking"})
	}

	return c.Status(fiber.StatusCreated).JSON(booking)
}


func UpdateBooking(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID tidak valid"})
	}

	var updated models.Booking
	if err := c.BodyParser(&updated); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format data salah"})
	}

	// Validasi status booking
	if !isValidBookingStatus(updated.Status) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Status booking tidak valid"})
	}

	updated.UpdatedAt = time.Now()

	update := bson.M{
		"$set": bson.M{
			"client_id":       updated.ClientID,
			"photographer_id": updated.PhotographerID,
			"date":            updated.Date,
			"location":        updated.Location,
			"status":          updated.Status,
			"note":            updated.Note,
			"updated_at":      updated.UpdatedAt,
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := bookingHandlerCollection.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal update booking"})
	}
	if result.MatchedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Booking tidak ditemukan"})
	}

	return c.JSON(fiber.Map{"message": "Booking diperbarui"})
}


func DeleteBooking(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID tidak valid"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := bookingHandlerCollection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal hapus booking"})
	}
	if result.DeletedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Booking tidak ditemukan"})
	}

	return c.JSON(fiber.Map{"message": "Booking dihapus"})
}
