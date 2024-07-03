package db

import (
	"context"

	"github.com/wellingtonchida/hotelreservation/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RoomStore interface {
	InsertRoom(ctx context.Context, room *types.Room) (*types.Room, error)
}

type MongoRoomStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoRoomStore(client *mongo.Client, db, coll string) *MongoRoomStore {
	return &MongoRoomStore{
		client: client,
		coll:   client.Database(db).Collection(coll),
	}
}

func (s *MongoRoomStore) InsertRoom(ctx context.Context, room *types.Room) (*types.Room, error) {
	resp, err := s.coll.InsertOne(ctx, room)
	if err != nil {
		return nil, err
	}
	room.ID = resp.InsertedID.(primitive.ObjectID)

	//TODO: update the hotel with this room id
	// filter := bson.M{"_id": room.HotelID}
	// update := bson.M{"$push": bson.M{"rooms": room.ID}}
	// _, err = s.client.Database("hotelreservation").Collection("hotels").UpdateOne(ctx, filter, update)

	return room, nil
}
