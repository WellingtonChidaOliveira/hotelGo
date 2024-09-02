package api

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/wellingtonchida/hotelreservation/types"
)

func TestPostUser(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	app := fiber.New()
	UserHandler := NewUserHandler(tdb.User)
	app.Post("/", UserHandler.HandlerPostUser)

	params := types.CreateUserRequest{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "jgh@jgh.com",
		Password:  "password",
	}
	b, _ := json.Marshal(params)
	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	var user types.User
	json.NewDecoder(resp.Body).Decode(&user)
	if len(user.ID) == 0 {
		t.Fatalf("expected id to be set")
	}
	if len(user.EncryptedPassword) > 0 {
		t.Fatalf("expected encrypted password to be empty")
	}
	if user.FirstName != params.FirstName {
		t.Fatalf("expected %s, got %s", params.FirstName, user.FirstName)
	}
	if user.LastName != params.LastName {
		t.Fatalf("expected %s, got %s", params.LastName, user.LastName)
	}
	if user.Email != params.Email {
		t.Fatalf("expected %s, got %s", params.Email, user.Email)
	}

}
