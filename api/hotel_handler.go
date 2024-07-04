package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wellingtonchida/hotelreservation/db"
	"github.com/wellingtonchida/hotelreservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HotelHandler struct {
	store *db.Store
}

func NewHotelHandler(store *db.Store) *HotelHandler {
	return &HotelHandler{store: store}
}

func (h *HotelHandler) HandleGetHotels(c *fiber.Ctx) error {
	hotels, err := h.store.Hotel.GetHotels(c.Context())
	if err != nil {
		return err
	}
	return c.JSON(hotels)
}

func (h *HotelHandler) HandleGetHotelRooms(c *fiber.Ctx) error {
	id := c.Params("id")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	filter := bson.M{"hotel_id": oid}
	rooms, err := h.store.Room.GetRooms(c.Context(), filter)
	if err != nil {
		return err
	}
	return c.JSON(rooms)
}

func (h *HotelHandler) HandleGetHotelByID(c *fiber.Ctx) error {
	id := c.Params("id")
	hotel, err := h.store.Hotel.GetHotelByID(c.Context(), id)
	if err != nil {
		return err
	}
	return c.JSON(hotel)
}

func (h *HotelHandler) HandlePostHotel(c *fiber.Ctx) error {
	var hotel types.Hotel
	if err := c.BodyParser(&hotel); err != nil {
		return err
	}
	insertedHotel, err := h.store.Hotel.InsertHotel(c.Context(), &hotel)
	if err != nil {
		return err
	}
	return c.JSON(insertedHotel)
}

func (h *HotelHandler) HandleUpdateHotel(c *fiber.Ctx) error {
	id := c.Params("id")
	var update bson.M
	if err := c.BodyParser(&update); err != nil {
		return err
	}
	if err := h.store.Hotel.UpdateHotel(c.Context(), bson.M{"_id": id}, update); err != nil {
		return err
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *HotelHandler) HandleDeleteHotel(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.store.Hotel.DeleteHotel(c.Context(), id); err != nil {
		return err
	}
	return c.SendStatus(fiber.StatusNoContent)
}
