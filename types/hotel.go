package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type Hotel struct {
	ID       primitive.ObjectID   `json:"_id,omitempty" bson:"_id,omitempty"`
	Name     string               `json:"name,omitempty" bson:"name,omitempty"`
	Location string               `json:"location,omitempty" bson:"location,omitempty"`
	Rooms    []primitive.ObjectID `json:"rooms,omitempty" bson:"rooms,omitempty"`
	Rating   int                  `json:"rating,omitempty" bson:"rating,omitempty"`
}
type RoomType int

const (
	_ RoomType = iota
	SingleRoomType
	DoubleRoomType
	SeaSideRoomType
	DeluxeRoomType
)

type Room struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Type      RoomType           `json:"type,omitempty" bson:"type,omitempty"`
	BasePrice float64            `json:"baseprice,omitempty" bson:"baseprice,omitempty"`
	Price     float64            `json:"price,omitempty" bson:"price,omitempty"`
	HotelID   primitive.ObjectID `json:"hotel_id,omitempty" bson:"hotel_id,omitempty"`
}
