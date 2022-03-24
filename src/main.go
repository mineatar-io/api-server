package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/buaazp/fasthttprouter"
	"github.com/joho/godotenv"
	"github.com/mineatar-io/api-server/src/conf"
	"github.com/mineatar-io/api-server/src/redis"
	"github.com/mineatar-io/api-server/src/routes"
	"github.com/mineatar-io/api-server/src/util"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/valyala/fasthttp"
)

var (
	host           string              = "127.0.0.1"
	port           uint16              = 3000
	prometheusHost string              = "127.0.0.1"
	prometheusPort uint16              = 3001
	config         *conf.Configuration = &conf.Configuration{}
	r              *redis.Redis        = &redis.Redis{}
)

var (
	requestCountMetric = promauto.NewCounter(prometheus.CounterOpts{
		Name: "request_count",
		Help: "The total number of requests received",
	})
)

func init() {
	var err error

	godotenv.Load()

	if err = config.ReadFile("config.yml"); err != nil {
		log.Fatal(err)
	}

	start := time.Now()

	if err = r.Connect(config.Redis.URI, config.Redis.Database); err != nil {
		log.Fatal(err)
	}

	log.Printf("Successfully connected to Redis (%s)\n", time.Since(start))

	if value, ok := os.LookupEnv("HOST"); ok {
		host = value
		prometheusHost = value
	}

	if value, ok := os.LookupEnv("PROM_HOST"); ok {
		prometheusHost = value
	}

	if value, ok := os.LookupEnv("PORT"); ok {
		parsedValue, err := strconv.ParseUint(value, 10, 16)

		if err != nil {
			log.Fatal(err)
		}

		port = uint16(parsedValue)
		prometheusPort = uint16(parsedValue + 1)
	}

	if value, ok := os.LookupEnv("PROM_PORT"); ok {
		parsedValue, err := strconv.ParseUint(value, 10, 16)

		if err != nil {
			log.Fatal(err)
		}

		prometheusPort = uint16(parsedValue)
	}

	routes.Init(r, config)
	util.Init(r, config)
}

func middleware(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
		ctx.Response.Header.Set("Access-Control-Allow-Headers", "*")
		ctx.Response.Header.Set("Access-Control-Allow-Methods", "GET,POST,HEAD,OPTIONS")
		ctx.Response.Header.Set("Access-Control-Expose-Headers", "X-Cache-Hit")

		if util.Debug {
			log.Printf("%s %s (%s) - %s\n", ctx.Method(), ctx.URI(), ctx.RemoteAddr(), ctx.UserAgent())
		}

		requestCountMetric.Inc()

		next(ctx)
	}
}

func metrics() {
	log.Printf("Prometheus listening on %s:%d\n", prometheusHost, prometheusPort)

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(fmt.Sprintf("%s:%d", prometheusHost, prometheusPort), nil)
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

	log.Printf("Listening on %s:%d\n", host, port)

	go metrics()

	log.Fatal(fasthttp.ListenAndServe(fmt.Sprintf("%s:%d", host, port), middleware(router.Handler)))
}
