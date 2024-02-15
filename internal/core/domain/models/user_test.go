package models

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
			ID:        id,
			FirstName: fn,
			LastName:  ln,
			Email:     email,
			Age:       age,
		}
		got := NewUser(id, fn, ln, email, age)
		err := got.Validate()
		if err != nil || !reflect.DeepEqual(wanted, got) {
			t.Errorf("failed to create valid user")
		}
	})

	t.Run("return error when empty fields is supplied", func(t *testing.T) {
		user := User{ID: "123"}
		err := user.Validate()
		if err == nil {
			t.Error(err)
		}
	})

	t.Run("return error when wrong uuid format is supplied", func(t *testing.T) {
		user := User{ID: "12345"}
		got := user.Validate().Error()
		wanted := "invalid id"
		if got != wanted {
			t.Errorf("got %v, wanted %v", got, wanted)
		}
	})

	t.Run("return error when wrong email format is supplied", func(t *testing.T) {
		id := uuid.NewString()
		user := User{
			ID:        id,
			FirstName: "first",
			LastName:  "last",
			Email:     "wexinc.com",
			Age:       20,
		}
		err := user.Validate()
		if err.Error() != "invalid email format" {
			t.Error(err)
		}
	})

	t.Run("return error when wrong email format is supplied", func(t *testing.T) {
		user := NewUser(uuid.NewString(), "john", "Doe", "john@ex.com", 17)
		got := user.Validate().Error()
		wanted := "user must be at least 18 years old"
		if got != wanted {
			t.Errorf("got %v, wanted %v", got, wanted)
		}
	})
}
