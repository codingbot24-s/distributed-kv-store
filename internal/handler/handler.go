package handler

import (
	"github.com/codingbot24-s/distributed-kv-store/internal/helper"
	"github.com/gofiber/fiber/v2"
)

type set struct {
	Key   string
	Value string
}

func Set(c *fiber.Ctx) error {
	var s set
	if err := c.BodyParser(&s); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":  err,
			"detail": "error parsing request body",
		})
	}

	cmd := helper.Command{
		OP:    "set",
		Key:   s.Key,
		Value: s.Value,
	}
	err := helper.ApplyCommand(&cmd)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":  err,
			"detail": "error applying command",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"key":     s.Key,
		"value":   s.Value,
	})
}
func Get(c *fiber.Ctx) error {
	key := c.Query("key")
	e, err := helper.GetEngine()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":  err,
			"detail": "error getting engine",
		})
	}
	res, ok := e.Get(key)
	if !ok {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":  err,
			"detail": "key not found",
		})
	}
	return c.JSON(fiber.Map{
		"success": true,
		"key":     key,
		"value":   res,
	})
}
func Delete() {}
