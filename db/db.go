package db

const (
	DBNAME     = "hotelreservation"
	TestDBNAME = "hotelreservation_test"
	DBURI      = "mongodb://localhost:27017"
)

type Store struct {
	User  UserStore
	Hotel HotelStore
	Room  RoomStore
}
