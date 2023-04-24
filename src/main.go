package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

var (
	app *fiber.App = fiber.New(fiber.Config{
		DisableStartupMessage: true,
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			log.Println(ctx.Request().URI(), err)

			return ctx.SendStatus(http.StatusInternalServerError)
		},
	})
	r    *Redis  = &Redis{}
	conf *Config = &Config{}
)

func init() {
	var err error

	if err = conf.ReadFile("config.yml"); err != nil {
		log.Fatal(err)
	}

	if err = r.Connect(conf.Redis); err != nil {
		log.Fatal(err)
	}

	log.Println("Successfully connected to Redis")

	app.Use(recover.New())

	if conf.Environment == "development" {
		app.Use(cors.New(cors.Config{
			AllowOrigins:  "*",
			AllowMethods:  "HEAD,OPTIONS,GET",
			ExposeHeaders: "X-Cache-Hit,X-Cache-Time-Remaining",
		}))

		app.Use(logger.New(logger.Config{
			Format:     "${time} ${ip}:${port} -> ${status}: ${method} ${path} (${latency})\n",
			TimeFormat: "2006/01/02 15:04:05",
		}))
	}
}

func main() {
	defer r.Close()

	instanceID, err := GetInstanceID()

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Listening on %s:%d\n", conf.Host, conf.Port+instanceID)
	log.Fatal(app.Listen(fmt.Sprintf("%s:%d", conf.Host, conf.Port+instanceID)))
}
