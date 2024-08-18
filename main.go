package main

import (
	"context"
	"flag"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/wellingtonchida/hotelreservation/api"
	"github.com/wellingtonchida/hotelreservation/api/middleware"
	"github.com/wellingtonchida/hotelreservation/db"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var config = fiber.Config{
	ErrorHandler: func(c *fiber.Ctx, err error) error {
		return c.JSON(map[string]string{"error": err.Error()})
	},
}

func main() {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}
	//handler initialization
	var (
		listenAddr = flag.String("listen", ":5000", "server listen address")
		userStore  = db.NewMongoUserStore(client, "users")
		hotelStore = db.NewMongoHotelStore(client, "hotels")
		roomStore  = db.NewMongoRoomStore(client, "rooms", hotelStore)
		store      = &db.Store{
			Hotel: hotelStore,
			Room:  roomStore,
			User:  userStore,
		}
		userHandler  = api.NewUserHandler(userStore)
		hotelHandler = api.NewHotelHandler(store)
		authHandler  = api.NewAuthHandler(userStore)
		roomHandler  = api.NewRoomHandler(store)
		app          = fiber.New(config)
		auth         = app.Group("/api")
		apiv1        = app.Group("/api/v1", middleware.JWTAuthentication(userStore))
	)
	flag.Parse()

	auth.Post("/auth", authHandler.HandleAuthenticate)

	apiv1.Get("/user", userHandler.HandleGetUsers)
	apiv1.Get("/user/:id", userHandler.HandleGetUser)
	apiv1.Post("/user", userHandler.HandlerPostUser)
	apiv1.Put("/user/:id", userHandler.HandlerPutUser)
	apiv1.Delete("/user/:id", userHandler.HandlerDeleteUser)

	apiv1.Get("/hotels", hotelHandler.HandleGetHotels)
	apiv1.Get("/hotels/:id", hotelHandler.HandleGetHotelByID)
	apiv1.Post("/hotels", hotelHandler.HandlePostHotel)
	apiv1.Put("/hotels/:id", hotelHandler.HandleUpdateHotel)
	apiv1.Delete("/hotels/:id", hotelHandler.HandleDeleteHotel)
	apiv1.Get("/hotels/:id/rooms", hotelHandler.HandleGetHotelRooms)

	apiv1.Post("/room/:id/book", roomHandler.HandleBookRoom)
	app.Listen(*listenAddr)
}
