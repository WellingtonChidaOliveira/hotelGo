package api

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/wellingtonchida/hotelreservation/db"
	"github.com/wellingtonchida/hotelreservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BookRoomParams struct {
	FromDate      time.Time `json:"fromDate"`
	TillDate      time.Time `json:"tillDate"`
	NumberPersons int       `json:"numberPersons"`
}

func (p BookRoomParams) Validate() error {
	now := time.Now()
	if p.FromDate.Before(now) {
		return fiber.NewError(fiber.StatusBadRequest, "fromDate must be in the future")
	}
	if p.TillDate.Before(p.FromDate) {
		return fiber.NewError(fiber.StatusBadRequest, "tillDate must be after fromDate")
	}
	if p.FromDate.IsZero() {
		return fiber.NewError(fiber.StatusBadRequest, "fromDate is required")
	}
	if p.TillDate.IsZero() {
		return fiber.NewError(fiber.StatusBadRequest, "tillDate is required")
	}
	if p.NumberPersons <= 0 {
		return fiber.NewError(fiber.StatusBadRequest, "numberPersons must be greater than 0")
	}
	return nil
}

type RoomHandler struct {
	store *db.Store
}

func NewRoomHandler(store *db.Store) *RoomHandler {
	return &RoomHandler{store: store}
}

func (h *RoomHandler) HandleGetRooms(c *fiber.Ctx) error {
	rooms, err := h.store.Room.GetRooms(c.Context(), bson.M{})
	if err != nil {
		return err
	}
	return c.JSON(rooms)
}

func (h *RoomHandler) HandleBookRoom(c *fiber.Ctx) error {
	var params BookRoomParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}

	if err := params.Validate(); err != nil {
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

	ok, err = h.isRoomAvailable(c.Context(), roomID, params)
	if err != nil {
		return err
	}

	fmt.Printf("room %s available: %t\n", roomID.Hex(), ok)

	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(genericResp{
			Type:    "error",
			Message: fmt.Sprintf("room %s already booked", roomID.Hex()),
		})
	}

	booking := &types.Booking{
		UserID:        user.ID,
		RoomID:        roomID,
		FromDate:      params.FromDate,
		TillDate:      params.TillDate,
		NumberPersons: params.NumberPersons,
	}

	inserted, err := h.store.Booking.InsertBooking(c.Context(), booking)
	if err != nil {
		return err
	}

	return c.JSON(inserted)
}

func (h *RoomHandler) isRoomAvailable(ctx context.Context, roomID primitive.ObjectID, params BookRoomParams) (bool, error) {

	where := bson.M{
		"roomID":   roomID,
		"fromDate": bson.M{"$gte": params.FromDate},
		"tillDate": bson.M{"$lte": params.TillDate},
	}

	bookings, err := h.store.Booking.GetBookings(ctx, where)

	fmt.Println(len(bookings))
	if err != nil {
		return false, err
	}
	ok := len(bookings) == 0
	return ok, nil
}
