package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
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

	godotenv.Load()

	if err = conf.ReadFile("config.yml"); err != nil {
		log.Fatal(err)
	}

	if err = r.Connect(conf.Redis); err != nil {
		log.Fatal(err)
	}

	log.Println("Successfully connected to Redis")

	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:  "*",
		AllowMethods:  "HEAD,OPTIONS,GET",
		ExposeHeaders: "Content-Type,Content-Disposition,X-Cache-Hit",
	}))
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
