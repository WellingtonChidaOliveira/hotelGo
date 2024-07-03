package main

import (
	"context"
	"fmt"
	"log"

	"github.com/wellingtonchida/hotelreservation/db"
	"github.com/wellingtonchida/hotelreservation/types"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	ctx := context.Background()
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}

	hotelStore := db.NewMongoHotelStore(client, db.DBNAME, "hotels")
	roomStore := db.NewMongoRoomStore(client, db.DBNAME, "rooms")
	hotel := types.Hotel{
		Name:     "Hotel California",
		Location: "California",
	}

	rooms := []types.Room{
		{
			Type:      types.SingleRoomType,
			BasePrice: 100.00,
		},
		{
			Type:      types.DoubleRoomType,
			BasePrice: 150.00,
		},
		{
			Type:      types.SeaSideRoomType,
			BasePrice: 200.00,
		},
		{
			Type:      types.DeluxeRoomType,
			BasePrice: 250.00,
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

	fmt.Printf("inserted hotel: %+v\n", insertedHotel)

	fmt.Println("seeding the database...")
}
