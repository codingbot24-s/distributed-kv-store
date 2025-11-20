package router

import (
	"fmt"

	"github.com/codingbot24-s/distributed-kv-store/internal/handler"
	"github.com/gofiber/fiber/v2"
)

func StartRouter() {
	app := fiber.New()
	app.Post("/set", handler.Set)

	fmt.Println("server is running on port 8080")
	app.Listen(":8080")
}
