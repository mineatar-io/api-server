package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

var (
	app *fiber.App = fiber.New(fiber.Config{
		DisableStartupMessage: true,
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			var fiberError *fiber.Error

			if errors.As(err, &fiberError) {
				return ctx.SendStatus(fiberError.Code)
			}

			log.Printf("Error: %v - URI: %s\n", err, ctx.Request().URI())

			return ctx.SendStatus(http.StatusInternalServerError)
		},
	})
	r          *Redis  = &Redis{}
	config     *Config = &Config{}
	instanceID uint16  = 0
)

func init() {
	var err error

	if err = config.ReadFile("config.yml"); err != nil {
		log.Fatal(err)
	}

	if err = r.Connect(config.Redis); err != nil {
		log.Fatal(err)
	}

	log.Println("Successfully connected to Redis")

	if instanceID, err = GetInstanceID(); err != nil {
		panic(err)
	}
}

func main() {
	defer r.Close()

	log.Printf("Listening on %s:%d\n", config.Host, config.Port+instanceID)

	if err := app.Listen(fmt.Sprintf("%s:%d", config.Host, config.Port+instanceID)); err != nil {
		panic(err)
	}
}
