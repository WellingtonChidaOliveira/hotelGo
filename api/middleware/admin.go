package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wellingtonchida/hotelreservation/types"
)

func AdminAuth(c *fiber.Ctx) error {
	user, ok := c.Context().UserValue("user").(*types.User)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "user not found")
	}
	if !user.IsAdmin {
		return fiber.NewError(fiber.StatusUnauthorized, "user is not an admin")
	}
	return c.Next()
}
