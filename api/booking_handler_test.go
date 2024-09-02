package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/wellingtonchida/hotelreservation/db/fixtures"
	"github.com/wellingtonchida/hotelreservation/types"
)

func TestAdminGetBookings(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	user := fixtures.AddUser(tdb.Store, "John", "Doe", true)
	hotel := fixtures.AddHotel(tdb.Store, "Hilton", "New York", 5, nil)
	room := fixtures.AddRoom(tdb.Store, "small", false, 100.00, hotel.ID)
	booking := fixtures.AddBooking(tdb.Store, user.ID, room.ID, time.Now(), time.Now().AddDate(0, 0, 1))

	app := fiber.New()
	AdminHandler := NewBookingHandler(tdb.Store)
	app.Get("/", AdminHandler.HandleGetBookings)

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", CreateTokenFromUser(user)))
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status code 200, got %d", resp.StatusCode)
	}

	var bookings []types.Booking
	if err := json.NewDecoder(resp.Body).Decode(&bookings); err != nil {
		t.Fatal(err)
	}

	if len(bookings) != 1 {
		t.Fatalf("expected 1 booking, got %d", len(bookings))
	}

	if bookings[0].ID != booking.ID {
		t.Fatalf("expected %s, got %s", booking.ID.Hex(), bookings[0].ID)
	}

	if bookings[0].RoomID != room.ID {
		t.Fatalf("expected %s, got %s", room.ID.Hex(), bookings[0].RoomID)
	}

	if bookings[0].UserID != user.ID {
		t.Fatalf("expected %s, got %s", user.ID.Hex(), bookings[0].UserID)
	}

}
