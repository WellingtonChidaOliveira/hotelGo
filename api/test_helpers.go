package api

import (
	"context"
	"testing"

	"github.com/wellingtonchida/hotelreservation/db"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	testdburi = "mongodb://localhost:27017"
)

type testdb struct {
	client *mongo.Client
	*db.Store
}

func (tdb *testdb) teardown(t *testing.T) {
	if err := tdb.client.Database(db.TestDBNAME).Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}
	//	if err := tdb.client.Disconnect(context.TODO()); err != nil {
	//		t.Fatal(err)
	//	}
}

func setup(t *testing.T) *testdb {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(testdburi))
	if err != nil {
		t.Fatal(err)
	}
	hotelStore := db.NewMongoHotelStore(client, db.TestDBNAME, "hotels")
	roomStore := db.NewMongoRoomStore(client, db.TestDBNAME, "rooms", hotelStore)
	return &testdb{
		client: client,
		Store: &db.Store{
			User:    db.NewMongoUserStore(client, db.TestDBNAME, "users"),
			Booking: db.NewMongoBookingStore(client, db.TestDBNAME, "bookings", roomStore),
			Hotel:   hotelStore,
			Room:    roomStore,
		},
	}
}
