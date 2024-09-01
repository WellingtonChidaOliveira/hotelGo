package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wellingtonchida/hotelreservation/types"
)

func getAuthuser(c *fiber.Ctx) (*types.User, error) {
	user, ok := c.Context().UserValue("user").(*types.User)
	if !ok {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "user is not authorized")
	}
	return user, nil
}
