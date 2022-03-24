package routes

import (
	"log"
	"net/http"

	"github.com/mineatar-io/api-server/src/util"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/valyala/fasthttp"
)

var (
	uuidLookupMetric = promauto.NewCounter(prometheus.CounterOpts{
		Name: "uuid_lookup_request_count",
		Help: "The amount of UUID lookup requests",
	})
)

func UUIDHandler(ctx *fasthttp.RequestCtx) {
	user := ctx.UserValue("user").(string)

	uuid, ok, err := util.LookupUUID(user)

	if err != nil {
		log.Println(err)

		ctx.SetStatusCode(http.StatusInternalServerError)
		ctx.SetBodyString(http.StatusText(http.StatusInternalServerError))

		return
	}

	uuidLookupMetric.Inc()

	if !ok {
		ctx.SetStatusCode(404)
		ctx.SetBodyString(http.StatusText(http.StatusNotFound))

		return
	}

	ctx.SetBodyString(uuid)
}
