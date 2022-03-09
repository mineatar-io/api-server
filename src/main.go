package main

import (
	"fmt"
	"log"
	"main/src/redis"
	"main/src/routes"
	"time"

	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

var (
	config *Configuration = &Configuration{}
	r      *redis.Redis   = &redis.Redis{}
)

func init() {
	var err error

	if err = config.ReadFile("config.yml"); err != nil {
		log.Fatal(err)
	}

	start := time.Now()

	if err = r.Connect(config.Redis.URI, config.Redis.Database); err != nil {
		log.Fatal(err)
	}

	log.Printf("Successfully connected to Redis (%s)\n", time.Since(start))

	routes.InitRoutes(r)
}

func main() {
	defer r.Close()

	router := fasthttprouter.New()

	router.GET("/ping", routes.PingHandler)
	router.GET("/uuid/:user", routes.UUIDHandler)
	router.GET("/skin/:user", routes.SkinHandler)
	router.GET("/face/:user", routes.FaceHandler)
	router.GET("/head/:user", routes.HeadHandler)
	router.GET("/body/full/:user", routes.FullBodyHandler)
	router.GET("/body/front/:user", routes.FrontBodyHandler)
	router.GET("/body/back/:user", routes.BackBodyHandler)
	router.GET("/body/left/:user", routes.LeftBodyHandler)
	router.GET("/body/right/:user", routes.RightBodyHandler)

	log.Printf("Listening on %s:%d\n", config.Host, config.Port)

	log.Fatal(fasthttp.ListenAndServe(fmt.Sprintf("%s:%d", config.Host, config.Port), router.Handler))
}
