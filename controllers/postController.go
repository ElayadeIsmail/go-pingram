package controllers

import (
	"fmt"

	"github.com/ElayadeIsmail/go-pingram/database"
	"github.com/ElayadeIsmail/go-pingram/models"
	"github.com/gofiber/fiber/v2"
)

func GetPosts(c *fiber.Ctx) error {
	var posts []models.Post
	result := database.DB.Joins("User").Find(&posts)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": result.Error,
			"data":    nil,
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "ok",
		"data": fiber.Map{
			"count": result.RowsAffected,
			"posts": posts,
		},
	})
}

func AddPost(c *fiber.Ctx) error {
	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
			"data":    nil,
		})
	}
	if found, ok := data["text"]; !ok || len(found) < 2 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Bad Request", "data": fiber.Map{"text": "Text must be more that 2 char"}})
	}

	userId, ok := c.Locals("userId").(int)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "something went wrong",
			"data":    nil,
		})
	}
	var u models.User
	if err := database.DB.First(&u, userId).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
			"data":    nil,
		})
	}
	fmt.Println()
	p := models.Post{
		Text:     data["text"],
		ImageUrl: data["imageUrl"],
		User:     u,
	}
	if err := database.DB.Create(&p).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
			"data":    nil,
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "ok",
		"data":    p,
	})
}

func DeletePost(c *fiber.Ctx) error {
	postId, err := c.ParamsInt("id")
	userId := c.Locals("userId").(int)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
			"data":    nil,
		})
	}
	var p models.Post
	if err := database.DB.Select("").Find(&p, postId).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
			"data":    nil,
		})
	}

	if p.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "Record Not Found",
			"data":    nil,
		})
	}

	if p.UserID != uint(userId) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Unauthorized for this action",
			"data":    nil,
		})
	}
	database.DB.Delete(&p)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Record was deleted successfully",
		"data":    p,
	})
}

func GetPostById(c *fiber.Ctx) error {
	postId, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid params type",
			"data":    nil,
		})
	}
	var p models.Post
	database.DB.Joins("User").Find(&p, postId)
	if p.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "Post not found",
			"data":    nil,
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Post was found",
		"data":    p,
	})
}
