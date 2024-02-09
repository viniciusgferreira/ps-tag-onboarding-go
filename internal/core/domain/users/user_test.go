package domain

import (
	"github.com/google/uuid"
	"reflect"
	"testing"
)

func TestUser(t *testing.T) {
	t.Run("create a valid user", func(t *testing.T) {
		const (
			id    = "20dc141d-3108-4529-8d65-ff3b046954be"
			fn    = "John"
			ln    = "Doe"
			email = "john@wex.com"
			age   = 31
		)
		wanted := User{
			id:        id,
			firstName: fn,
			lastName:  ln,
			email:     email,
			age:       age,
		}
		got := NewUser(id, fn, ln, email, age)
		err := got.validate()
		if err != nil || !reflect.DeepEqual(wanted, got) {
			t.Errorf("failed to create valid user")
		}
	})

	t.Run("return error when empty fields is supplied", func(t *testing.T) {
		user := User{id: "123"}
		err := user.validate()
		if err == nil {
			t.Error(err)
		}
	})

	t.Run("return error when wrong uuid format is supplied", func(t *testing.T) {
		user := User{id: "12345"}
		got := user.validate().Error()
		wanted := "invalid id"
		if got != wanted {
			t.Errorf("got %v, wanted %v", got, wanted)
		}
	})

	t.Run("return error when wrong email format is supplied", func(t *testing.T) {
		id := uuid.NewString()
		user := User{
			id:        id,
			firstName: "first",
			lastName:  "last",
			email:     "wexinc.com",
			age:       20,
		}
		err := user.validate()
		if err.Error() != "invalid email format" {
			t.Error(err)
		}
	})

	t.Run("return error when wrong email format is supplied", func(t *testing.T) {
		user := NewUser(uuid.NewString(), "john", "Doe", "john@ex.com", 17)
		got := user.validate().Error()
		wanted := "user must be at least 18 years old"
		if got != wanted {
			t.Errorf("got %v, wanted %v", got, wanted)
		}
	})
}
