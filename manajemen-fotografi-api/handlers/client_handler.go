package handlers

import (
	"context"
	
	"time"

	"manajemen-fotografi-api/config"
	"manajemen-fotografi-api/models"
	"manajemen-fotografi-api/utils"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var clientCollection = config.GetCollection("clients")

func CreateClient(c *fiber.Ctx) error {
    var client models.Client

    if err := c.BodyParser(&client); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format data salah"})
    }

    if client.UserID.IsZero() {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "User ID harus diisi"})
    }

    if len(client.Name) < 3 {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Nama client minimal 3 karakter"})
    }

    if !utils.IsValidPhone(client.Phone) {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Nomor telepon tidak valid"})
    }

    if len(client.Address) < 5 {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Alamat minimal 5 karakter"})
    }

    client.ID = primitive.NewObjectID()
    client.CreatedAt = time.Now().Unix()
    client.UpdatedAt = client.CreatedAt

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    _, err := clientCollection.InsertOne(ctx, client)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal menyimpan client"})
    }

    return c.Status(fiber.StatusCreated).JSON(client)
}





// GetClientByID mendapatkan data client berdasarkan ID
func GetClientByID(c *fiber.Ctx) error {
	idParam := c.Params("id")
	clientID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID tidak valid"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var client models.Client
	err = clientCollection.FindOne(ctx, bson.M{"_id": clientID}).Decode(&client)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Client tidak ditemukan"})
	}

	return c.JSON(client)
}

// GetClientByUserID mendapatkan data client berdasarkan UserID (relasi)
func GetClientByUserID(c *fiber.Ctx) error {
	userIDParam := c.Params("user_id")
	userID, err := primitive.ObjectIDFromHex(userIDParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "User ID tidak valid"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var client models.Client
	err = clientCollection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&client)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Client tidak ditemukan"})
	}

	return c.JSON(client)
}

// GetAllClients mengambil semua data client
func GetAllClients(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := clientCollection.Find(ctx, bson.M{})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal mengambil data client"})
	}
	defer cursor.Close(ctx)

	var clients []models.Client
	if err = cursor.All(ctx, &clients); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal memproses data client"})
	}

	return c.JSON(clients)
}

// UpdateClient mengubah data client berdasarkan ID
func UpdateClient(c *fiber.Ctx) error {
    idParam := c.Params("id")
    clientID, err := primitive.ObjectIDFromHex(idParam)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID tidak valid"})
    }

    var updateData models.Client
    if err := c.BodyParser(&updateData); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format data salah"})
    }

    // Validasi Name jika ada update
    if updateData.Name != "" && len(updateData.Name) < 3 {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Nama client minimal 3 karakter"})
    }

    // Validasi Phone jika ada update
    if updateData.Phone != "" && !utils.IsValidPhone(updateData.Phone) {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Nomor telepon tidak valid"})
    }

    // Validasi Address jika ada update
    if updateData.Address != "" && len(updateData.Address) < 5 {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Alamat minimal 5 karakter"})
    }

    updateData.UpdatedAt = time.Now().Unix()

    update := bson.M{
        "name":       updateData.Name,
        "phone":      updateData.Phone,
        "address":    updateData.Address,
        "updated_at": updateData.UpdatedAt,
    }

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    result, err := clientCollection.UpdateOne(ctx, bson.M{"_id": clientID}, bson.M{"$set": update})
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal update client"})
    }

    if result.MatchedCount == 0 {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Client tidak ditemukan"})
    }

    var updatedClient models.Client
    err = clientCollection.FindOne(ctx, bson.M{"_id": clientID}).Decode(&updatedClient)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal mengambil data yang diupdate"})
    }

    return c.JSON(updatedClient)
}


// DeleteClient menghapus data client berdasarkan ID
func DeleteClient(c *fiber.Ctx) error {
	idParam := c.Params("id")
	clientID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID tidak valid"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := clientCollection.DeleteOne(ctx, bson.M{"_id": clientID})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal menghapus client"})
	}

	if result.DeletedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Client tidak ditemukan"})
	}

	return c.JSON(fiber.Map{"message": "Client berhasil dihapus"})
}