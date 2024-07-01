package main

import (
	"flag"

	"github.com/gofiber/fiber/v2"
	"github.com/wellingtonchida/hotelreservation/api"
)

func main() {
	listenAddr := flag.String("listen", ":5000", "server listen address")
	flag.Parse()

	app := fiber.New()
	apiv1 := app.Group("/api/v1")

	apiv1.Get("/user", api.HandleGetUsers)    //list users
	apiv1.Get("/user/:id", api.HandleGetUser) //list user by id
	app.Listen(*listenAddr)
}
