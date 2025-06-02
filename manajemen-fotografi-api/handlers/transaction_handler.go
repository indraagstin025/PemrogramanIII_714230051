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



var (
	transactionCollection = config.GetTransactionCollection()
	bookingCollection     = config.GetBookingCollection()
)


func CreateDummyTransaction(c *fiber.Ctx) error {
	var trx models.Transaction

	if err := c.BodyParser(&trx); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format data salah"})
	}

	// Validasi data transaksi wajib
	if trx.BookingID.IsZero() || trx.Method == "" || trx.Total <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Data transaksi tidak lengkap"})
	}

	// Validasi metode pembayaran
	validMethods := map[string]bool{
		"transfer": true,
		"ewallet":  true,
	}

	if !validMethods[trx.Method] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Metode pembayaran tidak valid"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Cek apakah booking tersedia dan statusnya pending
	var booking models.Booking
	err := bookingCollection.FindOne(ctx, bson.M{"_id": trx.BookingID}).Decode(&booking)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Booking tidak ditemukan"})
	}

	if booking.Status != models.BookingStatusPending {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Booking tidak dalam status pending"})
	}

	now := time.Now()
	trx.ID = primitive.NewObjectID()
	trx.Status = "paid"
	trx.CreatedAt = now
	trx.UpdatedAt = now

	_, err = transactionCollection.InsertOne(ctx, trx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal menyimpan transaksi"})
	}

	// Update status booking menjadi confirmed
	updateResult, err := bookingCollection.UpdateOne(ctx, bson.M{"_id": trx.BookingID}, bson.M{
		"$set": bson.M{
			"status":     models.BookingStatusConfirmed,
			"updated_at": now,
		},
	})
	if err != nil {
		// Warning: transaksi sudah tersimpan tapi update booking gagal
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message":     "Transaksi berhasil, tapi gagal update status booking",
			"transaction": trx,
		})
	}

	if updateResult.MatchedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Booking tidak ditemukan saat update status"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":     "Transaksi berhasil (dummy), booking dikonfirmasi",
		"transaction": trx,
	})
}

func GetAllTransactions(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := transactionCollection.Find(ctx, bson.M{})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal mengambil data"})
	}
	defer cursor.Close(ctx)

	var transactions []models.Transaction
	if err := cursor.All(ctx, &transactions); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal decode data"})
	}

	return c.JSON(transactions)
}

// Contoh tambahan: dapatkan transaksi berdasarkan ID
func GetTransactionByID(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID tidak valid"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var trx models.Transaction
	err = transactionCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&trx)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Transaksi tidak ditemukan"})
	}

	return c.JSON(trx)
}
