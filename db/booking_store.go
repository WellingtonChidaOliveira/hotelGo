package db

import (
	"context"

	"github.com/wellingtonchida/hotelreservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BookingStore interface {
	InsertBooking(ctx context.Context, booking *types.Booking) (*types.Booking, error)
	GetBookings(ctx context.Context, filter bson.M) ([]*types.Booking, error)
	// GetBookingByID(ctx context.Context, id primitive.ObjectID) (*types.Booking, error)
	// GetBookingsByUserID(ctx context.Context, userID primitive.ObjectID) ([]*types.Booking, error)
	// GetBookingsByRoomID(ctx context.Context, roomID primitive.ObjectID) ([]*types.Booking, error)
}

type MongoBookingStore struct {
	client *mongo.Client
	coll   *mongo.Collection

	RoomStore
}

func NewMongoBookingStore(client *mongo.Client, coll string, roomStore RoomStore) *MongoBookingStore {
	return &MongoBookingStore{
		client:    client,
		coll:      client.Database(DBNAME).Collection(coll),
		RoomStore: roomStore,
	}
}

func (s *MongoBookingStore) InsertBooking(ctx context.Context, booking *types.Booking) (*types.Booking, error) {
	resp, err := s.coll.InsertOne(ctx, booking)
	if err != nil {
		return nil, err
	}
	booking.ID = resp.InsertedID.(primitive.ObjectID)

	return booking, nil
}

func (s *MongoBookingStore) GetBookings(ctx context.Context, filter bson.M) ([]*types.Booking, error) {
	cursor, err := s.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var bookings []*types.Booking
	if err := cursor.All(ctx, &bookings); err != nil {
		return nil, err
	}

	return bookings, nil
}
