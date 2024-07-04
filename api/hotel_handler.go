package api

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/wellingtonchida/hotelreservation/db"
	"github.com/wellingtonchida/hotelreservation/types"
	"go.mongodb.org/mongo-driver/bson"
)

type HotelHandler struct {
	hotelStore db.HotelStore
	roomStore  db.RoomStore
}

func NewHotelHandler(hs db.HotelStore, rs db.RoomStore) *HotelHandler {
	return &HotelHandler{hotelStore: hs, roomStore: rs}
}

type HotelQueryParams struct {
	Rooms    bool   `json:"rooms"`
	Rating   int    `json:"rating"`
	Location string `json:"location"`
}

func (h *HotelHandler) HandleGetHotels(c *fiber.Ctx) error {
	var qparams = HotelQueryParams{}
	if err := c.QueryParser(&qparams); err != nil {
		return err
	}
	log.Println(qparams)
	hotels, err := h.hotelStore.GetHotels(c.Context())
	if err != nil {
		return err
	}
	return c.JSON(hotels)
}

func (h *HotelHandler) HandleGetHotelByID(c *fiber.Ctx) error {
	id := c.Params("id")
	hotel, err := h.hotelStore.GetHotelByID(c.Context(), id)
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
	insertedHotel, err := h.hotelStore.InsertHotel(c.Context(), &hotel)
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
	if err := h.hotelStore.UpdateHotel(c.Context(), bson.M{"_id": id}, update); err != nil {
		return err
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *HotelHandler) HandleDeleteHotel(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.hotelStore.DeleteHotel(c.Context(), id); err != nil {
		return err
	}
	return c.SendStatus(fiber.StatusNoContent)
}
