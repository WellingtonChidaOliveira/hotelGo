package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/wellingtonchida/hotelreservation/api"
	"github.com/wellingtonchida/hotelreservation/db"
	"github.com/wellingtonchida/hotelreservation/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client       *mongo.Client
	roomStore    db.RoomStore
	hotelStore   db.HotelStore
	userStore    db.UserStore
	bookingStore db.BookingStore
	ctx          = context.Background()
)

func seedUser(fname, lname, email, password string, isAdmin bool) *types.User {
	user, err := types.CreateUserRequestToUser(&types.CreateUserRequest{
		FirstName: fname,
		LastName:  lname,
		Email:     email,
		Password:  password,
	})

	user.IsAdmin = isAdmin
	if err != nil {
		log.Fatal(err)
	}

	insertedUser, err := userStore.InsertUser(ctx, user)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s -> %+v\n", insertedUser.Email, api.CreateTokenFromUser(insertedUser))
	return insertedUser
}

func seedHotel(name, location string, rating int) *types.Hotel {
	hotel := types.Hotel{
		Name:     name,
		Location: location,
		Rooms:    []primitive.ObjectID{},
		Rating:   rating,
	}

	insertedHotel, err := hotelStore.InsertHotel(ctx, &hotel)
	if err != nil {
		log.Fatal(err)
	}

	return insertedHotel

}

func seedRoom(hotelID primitive.ObjectID, size string, price float64, seaSide bool) *types.Room {
	room := types.Room{
		HotelID: hotelID,
		SeaSide: seaSide,
		Size:    size,
		Price:   price,
	}

	insertedRoom, err := roomStore.InsertRoom(ctx, &room)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("inserted room: %+v\n", insertedRoom.ID.Hex())

	return insertedRoom
}

func seedBooking(userID, roomID primitive.ObjectID, checkIn, checkOut time.Time) *types.Booking {

	booking := types.Booking{
		UserID:   userID,
		RoomID:   roomID,
		FromDate: checkIn,
		TillDate: checkOut,
	}
	if bookingStore == nil {
		fmt.Println("booking store is nil")
		bookingStore = db.NewMongoBookingStore(client, "bookings", roomStore)
	}
	insertedBooking, err := bookingStore.InsertBooking(ctx, &booking)
	if err != nil {
		log.Fatal(err)
	}
	return insertedBooking
}

func main() {
	user := seedUser("james", "bond", "jd@jd.com", "password", false)
	seedUser("admin", "admin", "admin@admin.com", "Admin@123", true)
	fmt.Printf("inserted user: %+v\n", user.ID.Hex())

	hotel := seedHotel("Hilton", "New York", 5)
	seedHotel("Marriot", "San Francisco", 4)
	seedHotel("Sheraton", "Los Angeles", 3)
	fmt.Printf("inserted hotel: %+v\n", hotel.ID.Hex())

	rentRoom := seedRoom(hotel.ID, "small", 100.00, false)
	seedRoom(hotel.ID, "medium", 200.00, false)
	seedRoom(hotel.ID, "large", 300.00, true)

	bookingRoom := seedBooking(user.ID, rentRoom.ID, time.Now(), time.Now().AddDate(0, 0, 2))
	fmt.Printf("inserted booking: %+v\n", bookingRoom.ID.Hex())

}

func init() {
	var err error
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}
	if err := client.Database(db.DBNAME).Drop(ctx); err != nil {
		log.Fatal(err)
	}
	hotelStore = db.NewMongoHotelStore(client, "hotels")
	roomStore = db.NewMongoRoomStore(client, "rooms", hotelStore)
	userStore = db.NewMongoUserStore(client, "users")
	bookingStore = db.NewMongoBookingStore(client, "bookings", roomStore)
}
