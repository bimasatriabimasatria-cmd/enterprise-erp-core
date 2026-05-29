package controllers

import (
	"os"
	"time"

	"enterprise-erp/config"
	"enterprise-erp/models"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// Struktur data yang diterima saat Register
type RegisterInput struct {
	TenantID string `json:"tenant_id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

func Register(c *fiber.Ctx) error {
	var input RegisterInput
	
	// Tangkap data dari permintaan pengguna
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Format data tidak valid"})
	}

	// Acak (Hash) Password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), 10)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal memproses password"})
	}

	// Simpan User Baru ke Database
	user := models.User{
		TenantID: input.TenantID,
		Name:     input.Name,
		Email:    input.Email,
		Password: string(hashedPassword),
		Role:     input.Role,
	}

	if err := config.DB.Create(&user).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal menyimpan user, mungkin email sudah terdaftar"})
	}

	return c.Status(201).JSON(fiber.Map{"message": "Registrasi berhasil", "user_id": user.ID})
}

// Struktur data yang diterima saat Login
type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Login(c *fiber.Ctx) error {
	var input LoginInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Format data tidak valid"})
	}

	// Cari user berdasarkan email
	var user models.User
	if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Email tidak ditemukan"})
	}

	// Cek kecocokan password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Password salah"})
	}

	// Buat Kunci JWT (JSON Web Token)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":   user.ID,
		"tenant_id": user.TenantID,
		"role":      user.Role,
		"exp":       time.Now().Add(time.Hour * 24).Unix(), // Token berlaku 24 jam
	})

	// Tanda tangani token dengan JWT_SECRET
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal membuat token akses"})
	}

	return c.JSON(fiber.Map{
		"message": "Login berhasil",
		"token":   tokenString,
	})
}