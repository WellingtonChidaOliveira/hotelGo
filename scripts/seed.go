package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/wellingtonchida/hotelreservation/api"
	"github.com/wellingtonchida/hotelreservation/db"
	"github.com/wellingtonchida/hotelreservation/db/fixtures"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}
	if err := client.Database(db.DBNAME).Drop(ctx); err != nil {
		log.Fatal(err)
	}
	hotelStore := db.NewMongoHotelStore(client, db.DBNAME, "hotels")
	roomStore := db.NewMongoRoomStore(client, db.DBNAME, "rooms", hotelStore)
	db := &db.Store{
		User:    db.NewMongoUserStore(client, db.DBNAME, "users"),
		Booking: db.NewMongoBookingStore(client, db.DBNAME, "bookings", roomStore),
		Hotel:   hotelStore,
		Room:    roomStore,
	}

	use := fixtures.AddUser(db, "james", "bond", false)
	fmt.Printf("%s -> %s\n", use.Email, api.CreateTokenFromUser(use))
	admin := fixtures.AddUser(db, "admin", "admin", true)
	fmt.Printf("%s -> %s\n", admin.Email, api.CreateTokenFromUser(admin))

	hotel := fixtures.AddHotel(db, "Hilton", "New York", 5, nil)
	room := fixtures.AddRoom(db, "small", false, 100.00, hotel.ID)
	booking := fixtures.AddBooking(db, use.ID, room.ID, time.Now(), time.Now().AddDate(0, 0, 1))
	fmt.Printf("inserted booking: %+v\n", booking.ID.Hex())

}
