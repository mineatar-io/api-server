package routes

import (
	"fmt"
	"log"
	"net/http"

	"github.com/mineatar-io/api-server/src/util"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/valyala/fasthttp"
)

var (
	requestRawSkinMetric = promauto.NewCounter(prometheus.CounterOpts{
		Name: "raw_skin_request_count",
		Help: "The amount of raw skin requests",
	})
)

func SkinHandler(ctx *fasthttp.RequestCtx) {
	user := ctx.UserValue("user").(string)

	opts := util.ParseQueryParams(ctx, config.Routes.RawSkin)

	uuid, ok, err := util.LookupUUID(user)

	if err != nil {
		log.Println(err)

		ctx.SetStatusCode(http.StatusInternalServerError)
		ctx.SetBodyString(http.StatusText(http.StatusInternalServerError))

		return
	}

	if !ok && !opts.Fallback {
		ctx.SetStatusCode(http.StatusNotFound)
		ctx.SetBodyString(http.StatusText(http.StatusNotFound))

		return
	}

	rawSkin, _, err := util.GetPlayerSkin(uuid)

	if err != nil {
		log.Println(err)

		ctx.SetStatusCode(http.StatusInternalServerError)
		ctx.SetBodyString(http.StatusText(http.StatusInternalServerError))

		return
	}

	data, err := util.EncodePNG(rawSkin)

	if err != nil {
		log.Println(err)

		ctx.SetStatusCode(http.StatusInternalServerError)
		ctx.SetBodyString(http.StatusText(http.StatusInternalServerError))

		return
	}

	if opts.Download {
		ctx.Response.Header.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.png"`, user))
	}

	requestRawSkinMetric.Inc()

	ctx.SetContentType("image/png")
	ctx.SetBody(data)
}
