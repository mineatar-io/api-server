package routes

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/valyala/fasthttp"
)

var (
	pingRequestMetric = promauto.NewCounter(prometheus.CounterOpts{
		Name: "ping_request_count",
		Help: "The amount of ping requests",
	})
)

func PingHandler(ctx *fasthttp.RequestCtx) {
	pingRequestMetric.Inc()

	ctx.SetBodyString("Pong!")
}
