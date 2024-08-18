package api

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/wellingtonchida/hotelreservation/db"
	"github.com/wellingtonchida/hotelreservation/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BookRoomParams struct {
	FromDate      time.Time `json:"fromDate"`
	TillDate      time.Time `json:"tillDate"`
	NumberPersons int       `json:"numberPersons"`
}

type RoomHandler struct {
	store *db.Store
}

func NewRoomHandler(store *db.Store) *RoomHandler {
	return &RoomHandler{store: store}
}

func (h *RoomHandler) HandleBookRoom(c *fiber.Ctx) error {
	var params BookRoomParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}

	roomID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return err
	}

	user, ok := c.Context().Value("user").(*types.User)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(genericResp{
			Type:    "error",
			Message: "internal server error",
		})
	}

	booking := &types.Booking{
		UserID:        user.ID,
		RoomID:        roomID,
		FromDate:      params.FromDate,
		TillDate:      params.TillDate,
		NumberPersons: params.NumberPersons,
	}

	fmt.Printf("%+v\n", booking)
	return nil
}
