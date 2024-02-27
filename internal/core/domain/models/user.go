package models

type User struct {
	ID        string `bson:"_id,omitempty" json:"id,omitempty"`
	FirstName string `bson:"firstName" json:"firstName"`
	LastName  string `bson:"lastName" json:"lastName"`
	Email     string `bson:"email" json:"email"`
	Age       int    `bson:"age" json:"age"`
}

func NewUser(id, firstName, lastName, email string, age int) User {
	user := User{
		ID:        id,
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Age:       age,
	}
	return user
}
