package types

type User struct {
	ID        string `bson:"_id" json:"id",omitempty` //for marshalling and unmarshalling
	FirstName string `bson:"firstName" json:"firstName"`
	LastName  string `bson:"lastName" json:"lastName"`
	Email     string `bson:"email" json:"email"`
}
