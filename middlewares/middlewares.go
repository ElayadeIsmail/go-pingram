package middlewares

import (
	"github.com/ElayadeIsmail/go-pingram/config"
	"github.com/ElayadeIsmail/go-pingram/database"
	"github.com/ElayadeIsmail/go-pingram/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)

func CurrentUser(c *fiber.Ctx) error {
	cookie := c.Cookies(config.Getenv("COOKIE_NAME"))
	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(config.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		c.Locals("userId", 0)
		return c.Next()
	}
	claims := token.Claims.(*jwt.StandardClaims)
	var u models.User
	database.DB.Where("id = ?", claims.Issuer).First(&u)
	c.Locals("userId", int(u.ID))
	return c.Next()
}

func RequireLogin(c *fiber.Ctx) error {
	userId := c.Locals("userId")
	if userId == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Unauthorized",
			"data":    nil,
		})
	} else {
		return c.Next()
	}
}
