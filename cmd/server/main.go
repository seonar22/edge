package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/packetframe/auth"
	"github.com/packetframe/database"
	"github.com/packetframe/servicerouter"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/packetframe/edge"
)

// API listen port
const apiPort = 5002

// Validator
var validate *validator.Validate

// resp returns a JSON error response
func resp(ctx *fiber.Ctx, code int, message string) error {
	return ctx.Status(code).JSON(&fiber.Map{"errors": strings.Split(message, "\n")})
}

func main() {
	// Fiber API server
	app := fiber.New(fiber.Config{DisableStartupMessage: true})
	validate = validator.New()

	// Register validators
	if err := edge.RegisterValidators(validate); err != nil {
		log.Fatal(err)
	}

	db, err := database.New("mongodb://localhost:27017", "packetframe")
	if err != nil {
		log.Fatal(err)
	}

	// Register service with servicerouter
	if err := servicerouter.Register(db, &servicerouter.Service{
		Path: "edge",
		Port: apiPort,
		Repo: "https://github.com/packetframe/edge",
	}); err != nil {
		log.Fatal(err)
	}

	// Add a container
	app.Post("/containers", func(ctx *fiber.Ctx) error {
		user, err := auth.GetUser(ctx, db)
		if err != nil {
			return resp(ctx, 403, "Unauthorized")
		}

		// Parse body into struct
		var container edge.Container
		if err := ctx.BodyParser(&container); err != nil {
			return resp(ctx, 400, err.Error())
		}

		// Validate zone struct
		err = validate.Struct(container)
		if err != nil {
			return resp(ctx, 400, err.Error())
		}

		container.ID = primitive.ObjectID{}             // Reset ObjectId
		container.Users = []primitive.ObjectID{user.ID} // Set users array to contain the container creator
		if container.Env == nil {
			container.Env = map[string]string{}
		}

		// Insert the new zone
		_, err = db.Db.Collection("containers").InsertOne(context.Background(), container)
		if err != nil {
			log.Warn(err)
			return resp(ctx, 500, "Unable to create container")
		}

		return ctx.Status(200).JSON(&fiber.Map{"message": "Container created"})
	})

	log.Println("Starting API")
	log.Fatal(app.Listen(fmt.Sprintf("127.0.0.1:%d", apiPort)))
}
