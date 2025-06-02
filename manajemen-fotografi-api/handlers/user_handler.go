package handlers

import (
	"context"
	"time"
	"regexp"
    "unicode"

	"manajemen-fotografi-api/config"
	"manajemen-fotografi-api/models"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

var userCollection = config.GetCollection("users")

// Fungsi validasi email menggunakan regex sederhana
func isValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

// Fungsi validasi nama (3-50 karakter, hanya huruf dan spasi)
func isValidName(name string) bool {
	if len(name) < 3 || len(name) > 50 {
		return false
	}
	for _, r := range name {
		if !(unicode.IsLetter(r) || unicode.IsSpace(r)) {
			return false
		}
	}
	return true
}

// RegisterUser handler untuk registrasi user baru dengan validasi input
func RegisterUser(c *fiber.Ctx) error {
	type registerInput struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}

	var input registerInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format data salah"})
	}

	// Validasi input
	if !isValidEmail(input.Email) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format email tidak valid"})
	}

	if len(input.Password) < 8 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Password minimal 8 karakter"})
	}

	if !isValidName(input.Name) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Nama harus 3-50 karakter dan hanya berisi huruf dan spasi"})
	}

	// Validasi role
	if input.Role != models.RoleClient && input.Role != models.RolePhotographer {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Role tidak valid"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// cek apakah email sudah terdaftar
	count, err := userCollection.CountDocuments(ctx, bson.M{"email": input.Email})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal cek email"})
	}
	if count > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Email sudah terdaftar"})
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal memproses password"})
	}

	user := models.User{
		ID:        primitive.NewObjectID(),
		Name:      input.Name,
		Email:     input.Email,
		Password:  string(hashedPassword),
		Role:      input.Role,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}

	_, err = userCollection.InsertOne(ctx, user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal menyimpan user"})
	}

	user.Password = ""
	return c.Status(fiber.StatusCreated).JSON(user)
}

// LoginUser handler untuk login dengan validasi input email dan password
func LoginUser(c *fiber.Ctx) error {
	type loginInput struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var input loginInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format data salah"})
	}

	// Validasi input login
	if !isValidEmail(input.Email) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format email tidak valid"})
	}

	if len(input.Password) < 8 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Password minimal 8 karakter"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user models.User
	err := userCollection.FindOne(ctx, bson.M{"email": input.Email}).Decode(&user)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Email atau password salah"})
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Email atau password salah"})
	}

	user.Password = ""

	return c.JSON(fiber.Map{
		"message": "Login berhasil",
		"user":    user,
	})
}


// LogoutUser handler untuk logout
func LogoutUser(c *fiber.Ctx) error {
	sess := c.Locals("session").(*session.Session)
	if err := sess.Destroy(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal logout"})
	}
	return c.JSON(fiber.Map{"message": "Logout berhasil"})
}
