package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wellingtonchida/hotelreservation/types"
)

func HandleGetUsers(c *fiber.Ctx) error {
	u := types.User{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "jd@jd.com",
	}
	return c.JSON(u)
}

func HandleGetUser(c *fiber.Ctx) error {
	return c.JSON(map[string]string{"user": "John Doe"})
}
