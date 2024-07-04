package db

import (
	"context"

	"github.com/wellingtonchida/hotelreservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type HotelStore interface {
	GetHotels(context.Context) ([]*types.Hotel, error)
	GetHotelByID(context.Context, string) (*types.Hotel, error)
	InsertHotel(context.Context, *types.Hotel) (*types.Hotel, error)
	UpdateHotel(context.Context, bson.M, bson.M) error
	DeleteHotel(context.Context, string) error
}

type MongoHotelSotre struct {
	cliente *mongo.Client
	coll    *mongo.Collection
}

func NewMongoHotelStore(client *mongo.Client, coll string) *MongoHotelSotre {
	return &MongoHotelSotre{
		cliente: client,
		coll:    client.Database(DBNAME).Collection(coll),
	}
}

func (s *MongoHotelSotre) GetHotels(ctx context.Context) ([]*types.Hotel, error) {
	cursor, err := s.coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var hotels []*types.Hotel
	if err := cursor.All(ctx, &hotels); err != nil {
		return nil, err
	}
	return hotels, nil
}

func (s *MongoHotelSotre) GetHotelByID(ctx context.Context, id string) (*types.Hotel, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var hotel types.Hotel
	if err := s.coll.FindOne(ctx, bson.M{"_id": objID}).Decode(&hotel); err != nil {
		return nil, err
	}
	return &hotel, nil
}

func (s *MongoHotelSotre) InsertHotel(ctx context.Context, hotel *types.Hotel) (*types.Hotel, error) {
	resp, err := s.coll.InsertOne(ctx, hotel)
	if err != nil {
		return nil, err
	}
	hotel.ID = resp.InsertedID.(primitive.ObjectID)
	return hotel, nil
}

func (s *MongoHotelSotre) UpdateHotel(ctx context.Context, filter, update bson.M) error {
	_, err := s.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (s *MongoHotelSotre) DeleteHotel(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = s.coll.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return err
	}
	return nil
}
