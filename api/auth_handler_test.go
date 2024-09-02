package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/wellingtonchida/hotelreservation/db/fixtures"
)

func TestAuthenticateSuccess(t *testing.T) {
	tbd := setup(t)
	defer tbd.teardown(t)
	insertedUser := fixtures.AddUser(tbd.Store, "Thomas", "Shelby", false)

	app := fiber.New()
	authHandler := NewAuthHandler(tbd.User)
	app.Post("/auth", authHandler.HandleAuthenticate)

	params := AuthParams{
		Email:    insertedUser.Email,
		Password: "Thomas_Shelby",
	}
	b, _ := json.Marshal(params)
	req := httptest.NewRequest("POST", "/auth", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status code 200, got %d", resp.StatusCode)
	}

	var authResp AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		t.Fatal(err)
	}

	if len(authResp.Token) == 0 {
		t.Fatalf("expected token to be set")
	}
	if authResp.Token == "" {
		t.Fatalf("expected token to be set")
	}

	//set the encrypted password to empty so we can compare the two structs
	//The password is not returned in the response
	insertedUser.EncryptedPassword = ""
	if !reflect.DeepEqual(insertedUser, authResp.User) {
		t.Fatalf("expected the user -> %+v \n to be the inserted user -> %+v", insertedUser, authResp.User)
	}
}

func TestAuthenticateWithWrongPasswordFailure(t *testing.T) {
	tbd := setup(t)
	defer tbd.teardown(t)
	fixtures.AddUser(tbd.Store, "John", "Doe", false)

	app := fiber.New()
	authHandler := NewAuthHandler(tbd.User)
	app.Post("/auth", authHandler.HandleAuthenticate)

	params := AuthParams{
		Email:    "Doe@John.com",
		Password: "passworddsds",
	}
	b, _ := json.Marshal(params)
	req := httptest.NewRequest("POST", "/auth", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected status code 401, got %d", resp.StatusCode)
	}

	var genericResp genericResp
	if err := json.NewDecoder(resp.Body).Decode(&genericResp); err != nil {
		t.Fatal(err)
	}

	if genericResp.Type != "error" {
		t.Fatalf("expected error type")
	}
	if genericResp.Message != "invalid credentials" {
		t.Fatalf("expected invalid credentials message")
	}

}
