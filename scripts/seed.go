package main

import (
	"context"
	"fmt"
	"log"

	"github.com/wellingtonchida/hotelreservation/db"
	"github.com/wellingtonchida/hotelreservation/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client     *mongo.Client
	roomStore  db.RoomStore
	hotelStore db.HotelStore
	ctx        = context.Background()
)

func SeedHotel(name, location string, rating int) {
	hotel := types.Hotel{
		Name:     name,
		Location: location,
		Rooms:    []primitive.ObjectID{},
		Rating:   rating,
	}

	rooms := []types.Room{
		{
			Size:  "small",
			Price: 100.00,
		},
		{
			Size:  "medium",
			Price: 150.00,
		},
		{
			Size:  "large",
			Price: 200.00,
		},
		{
			Size:  "large",
			Price: 250.00,
		},
	}

	insertedHotel, err := hotelStore.InsertHotel(ctx, &hotel)
	if err != nil {
		log.Fatal(err)
	}
	for _, room := range rooms {
		room.HotelID = insertedHotel.ID
		insertRoom, err := roomStore.InsertRoom(ctx, &room)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("inserted room: %+v\n", insertRoom)
	}

}

func main() {

	SeedHotel("Hilton", "New York", 5)
	SeedHotel("Marriot", "San Francisco", 4)
	SeedHotel("Sheraton", "Los Angeles", 3)
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
}
