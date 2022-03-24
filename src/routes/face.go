package routes

import (
	"fmt"
	"log"
	"net/http"

	"github.com/mineatar-io/api-server/src/util"
	"github.com/mineatar-io/skin-render"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/valyala/fasthttp"
)

var (
	renderFaceMetric = promauto.NewCounter(prometheus.CounterOpts{
		Name: "face_render_count",
		Help: "The amount of face renders",
	})
)

func FaceHandler(ctx *fasthttp.RequestCtx) {
	user := ctx.UserValue("user").(string)

	opts := util.ParseQueryParams(ctx, config.Routes.Face)

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

	cacheKey := fmt.Sprintf("result:face-%d-%t-%s", opts.Scale, opts.Overlay, uuid)

	{
		cache, ok, err := r.GetBytes(cacheKey)

		if err != nil {
			log.Println(err)

			ctx.SetStatusCode(http.StatusInternalServerError)
			ctx.SetBodyString(http.StatusText(http.StatusInternalServerError))

			return
		}

		if ok {
			if opts.Download {
				ctx.Response.Header.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.png"`, user))
			}

			if util.Debug {
				log.Printf("Retrieved cache for face render for '%s'\n", uuid)
			}

			ctx.Response.Header.Set("X-Cache-Hit", "TRUE")
			ctx.SetContentType("image/png")
			ctx.SetBody(cache)

			return
		}
	}

	rawSkin, slim, err := util.GetPlayerSkin(uuid)

	if err != nil {
		log.Println(err)

		ctx.SetStatusCode(http.StatusInternalServerError)
		ctx.SetBodyString(http.StatusText(http.StatusInternalServerError))

		return
	}

	render := skin.RenderFace(rawSkin, skin.Options{
		Overlay: opts.Overlay,
		Slim:    slim,
		Scale:   opts.Scale,
	})

	if util.Debug {
		log.Printf("Rendered face image for '%s'\n", uuid)
	}

	renderFaceMetric.Inc()

	data, err := util.EncodePNG(render)

	if err != nil {
		log.Println(err)

		ctx.SetStatusCode(http.StatusInternalServerError)
		ctx.SetBodyString(http.StatusText(http.StatusInternalServerError))

		return
	}

	if err = r.Set(cacheKey, data, config.Cache.RenderCacheDuration); err != nil {
		log.Println(err)

		ctx.SetStatusCode(http.StatusInternalServerError)
		ctx.SetBodyString(http.StatusText(http.StatusInternalServerError))

		return
	}

	if opts.Download {
		ctx.Response.Header.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.png"`, user))
	}

	ctx.Response.Header.Set("X-Cache-Hit", "FALSE")
	ctx.SetContentType("image/png")
	ctx.SetBody(data)
}
