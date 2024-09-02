package fixtures

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/wellingtonchida/hotelreservation/db"
	"github.com/wellingtonchida/hotelreservation/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddUser(store *db.Store, fname, lname string, isAdmin bool) *types.User {
	user, err := types.CreateUserRequestToUser(&types.CreateUserRequest{
		FirstName: fname,
		LastName:  lname,
		Email:     fmt.Sprintf("%s@%s.com", lname, fname),
		Password:  fmt.Sprintf("%s_%s", fname, lname),
	})

	user.IsAdmin = isAdmin
	if err != nil {
		log.Fatal(err)
	}

	insertedUser, err := store.User.InsertUser(context.Background(), user)
	if err != nil {
		log.Fatal(err)
	}

	return insertedUser
}

func AddHotel(store *db.Store, name, location string, rating int, rooms []primitive.ObjectID) *types.Hotel {
	var roomIds = rooms
	if len(rooms) == 0 {
		roomIds = []primitive.ObjectID{}
	}
	hotel := types.Hotel{
		Name:     name,
		Location: location,
		Rooms:    roomIds,
		Rating:   rating,
	}

	insertedHotel, err := store.Hotel.InsertHotel(context.TODO(), &hotel)
	if err != nil {
		log.Fatal(err)
	}
	return insertedHotel
}

func AddRoom(store *db.Store, size string, seaSide bool, price float64, hotelID primitive.ObjectID) *types.Room {
	room := types.Room{
		SeaSide: seaSide,
		Size:    size,
		Price:   price,
		HotelID: hotelID,
	}

	insertedRoom, err := store.Room.InsertRoom(context.Background(), &room)
	if err != nil {
		log.Fatal(err)
	}
	return insertedRoom
}

func AddBooking(store *db.Store, userID, roomID primitive.ObjectID, checkIn, checkOut time.Time) *types.Booking {
	booking := types.Booking{
		UserID:   userID,
		RoomID:   roomID,
		FromDate: checkIn,
		TillDate: checkOut,
	}

	insertedBooking, err := store.Booking.InsertBooking(context.Background(), &booking)
	if err != nil {
		log.Fatal(err)
	}
	return insertedBooking
}
