package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wellingtonchida/hotelreservation/db"
	"go.mongodb.org/mongo-driver/bson"
)

type BookingHandler struct {
	store *db.Store
}

func NewBookingHandler(store *db.Store) *BookingHandler {
	return &BookingHandler{store: store}
}

// TODO: this needs to be admin authorized
func (b *BookingHandler) HandleGetBookings(c *fiber.Ctx) error {
	bookings, err := b.store.Booking.GetBookings(c.Context(), bson.M{})
	if err != nil {
		return err
	}
	return c.JSON(bookings)
}

// TODO: this needs to be user authorized
func (b *BookingHandler) HandleGetBooking(c *fiber.Ctx) error {
	id := c.Params("id")
	booking, err := b.store.Booking.GetBookingByID(c.Context(), id)
	if err != nil {
		return err
	}
	user, err := getAuthuser(c)
	if err != nil {
		return err
	}
	if booking.UserID != user.ID {
		return fiber.NewError(fiber.StatusUnauthorized, "user is not authorized to view this booking")
	}

	return c.JSON(booking)
}

func (b *BookingHandler) HandleCancelBooking(c *fiber.Ctx) error {
	id := c.Params("id")
	booking, err := b.store.Booking.GetBookingByID(c.Context(), id)
	if err != nil {
		return err
	}
	user, err := getAuthuser(c)
	if err != nil {
		return err
	}
	if booking.UserID != user.ID {
		return fiber.NewError(fiber.StatusUnauthorized, "user is not authorized to cancel this booking")
	}
	if booking.Canceled {
		return fiber.NewError(fiber.StatusBadRequest, "booking already canceled")
	}
	booking.Canceled = true

	if err := b.store.Booking.UpdateBooking(c.Context(), booking.ID.Hex(), bson.M{"canceled": true}); err != nil {
		return err
	}
	return c.SendStatus(fiber.StatusNoContent)
}
