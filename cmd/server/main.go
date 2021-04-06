package main

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/packetframe/database"
	"github.com/packetframe/servicerouter"
	log "github.com/sirupsen/logrus"
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

	log.Println("Starting API")
	log.Fatal(app.Listen(fmt.Sprintf("127.0.0.1:%d", apiPort)))
}
