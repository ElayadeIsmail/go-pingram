package controllers

import (
	"errors"
	"strconv"
	"time"

	"github.com/ElayadeIsmail/go-pingram/config"
	"github.com/ElayadeIsmail/go-pingram/database"
	"github.com/ElayadeIsmail/go-pingram/models"
	"github.com/ElayadeIsmail/go-pingram/validation"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

func hashPassword(p string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(p), 14)
	return string(bytes), err
}

func compareHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func getUserByEmail(e string) (*models.User, error) {
	db := database.DB
	var user models.User
	if err := db.Where(&models.User{Email: e}).Find(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func getUserByUsername(u string) (*models.User, error) {
	db := database.DB
	var user models.User
	if err := db.Where(&models.User{Username: u}).Find(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func Register(c *fiber.Ctx) error {
	var data map[string]string
	db := database.DB
	if err := c.BodyParser(&data); err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": err.Error(), "data": nil})
	}
	if errs := validation.Validate(models.User{
		Username: data["username"],
		Email:    data["email"],
		Password: data["password"],
	}); len(errs) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "invalid inputs",
			"data":    errs,
		})
	}
	if u, err := getUserByEmail(data["email"]); u != nil || err != nil {
		if u.ID != 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Email Already Exist", "data": nil})
		} else if err != nil {
			return c.Status(500).JSON(fiber.Map{"status": "error", "message": err.Error(), "data": nil})
		}
	}
	if u, err := getUserByUsername(data["username"]); u != nil || err != nil {
		if u.ID != 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Username Already Exist", "data": nil})
		} else if err != nil {
			return c.Status(500).JSON(fiber.Map{"status": "error", "message": err.Error(), "data": nil})
		}
	}
	hash, err := hashPassword(data["password"])
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": err.Error(), "data": nil})
	}
	u := models.User{
		Username: data["username"],
		Email:    data["email"],
		Password: hash,
	}
	if err := db.Create(&u).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": err.Error(), "data": nil})
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    strconv.Itoa(int(u.ID)),
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	})

	token, err := claims.SignedString([]byte(config.Getenv("JWT_SECRET")))
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"status":  "error",
			"message": "Could Not Sign user",
			"data":    nil,
		})
	}
	cookie := fiber.Cookie{
		Name:     config.Getenv("COOKIE_NAME"),
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}
	c.Cookie(&cookie)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Logged in successfully", "data": u,
	})
}

func Login(c *fiber.Ctx) error {
	type loginField struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=6,max=22"`
	}
	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Could Not Get Data",
			"data":    nil,
		})
	}
	if errors := validation.Validate(loginField{
		Email:    data["email"],
		Password: data["password"],
	}); len(errors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Bad Request",
			"data":    errors,
		})
	}
	u, err := getUserByEmail(data["email"])
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error(), "data": nil})
	}
	if u.ID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Invalid Credentails", "data": nil})
	}
	if isMatch := compareHash(data["password"], u.Password); !isMatch {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Invalid Credentails", "data": nil})
	}
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    strconv.Itoa(int(u.ID)),
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	})

	token, err := claims.SignedString([]byte(config.Getenv("JWT_SECRET")))
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
			"data":    nil,
		})
	}
	cookie := fiber.Cookie{
		Name:     config.Getenv("COOKIE_NAME"),
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}
	c.Cookie(&cookie)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Logged in successfully", "data": u,
	})
}

func Logout(c *fiber.Ctx) error {
	cookie := fiber.Cookie{
		Name:     config.Getenv("COOKIE_NAME"),
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}
	c.Cookie(&cookie)
	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "user Logged out successfully", "data": nil})
}

func CurrentUser(c *fiber.Ctx) error {
	userId := c.Locals("userId")
	if userId == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Unauthorized",
			"data":    nil,
		})
	}
	var u models.User
	database.DB.Where("id = ?", userId).First(&u)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "user found",
		"data":    u,
	})
}
