package routes

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/mineatar-io/api-server/src/util"
	"github.com/mineatar-io/api-server/src/util/renders"
	"github.com/valyala/fasthttp"
)

func FaceHandler(ctx *fasthttp.RequestCtx) {
	user := ctx.UserValue("user").(string)

	download := ctx.QueryArgs().GetBool("download")

	scale, err := ctx.QueryArgs().GetUint("scale")

	if err != nil {
		scale = config.Routes.Face.DefaultScale
	}

	scale = util.Clamp(scale, config.Routes.Face.MinScale, config.Routes.Face.MaxScale)

	overlay := true

	if ctx.QueryArgs().Has("overlay") {
		overlay = ctx.QueryArgs().GetBool("overlay")
	}

	uuid, err := util.GetUUID(user)

	if err != nil {
		log.Println(err)

		ctx.SetStatusCode(http.StatusInternalServerError)
		ctx.SetBodyString(http.StatusText(http.StatusInternalServerError))

		return
	}

	cacheKey := fmt.Sprintf("result:face-%d-%t-%s", scale, overlay, uuid)

	cache, ok, err := r.GetBytes(cacheKey)

	if err != nil {
		log.Println(err)

		ctx.SetStatusCode(http.StatusInternalServerError)
		ctx.SetBodyString(http.StatusText(http.StatusInternalServerError))

		return
	}

	if ok {
		if download {
			ctx.Response.Header.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.png"`, user))
		}

		ctx.Response.Header.Set("X-Cache-Hit", "TRUE")
		ctx.SetContentType("image/png")
		ctx.SetBody(cache)

		return
	}

	skin, slim, err := util.GetPlayerSkin(uuid)

	if err != nil {
		log.Println(err)

		ctx.SetStatusCode(http.StatusInternalServerError)
		ctx.SetBodyString(http.StatusText(http.StatusInternalServerError))

		return
	}

	if skin == nil {
		skin = util.GetDefaultSkin(slim)
	}

	render := renders.RenderFace(skin, renders.RenderOptions{
		Overlay: overlay,
		Slim:    slim,
		Scale:   scale,
	})

	data, err := util.EncodePNG(render)

	if err != nil {
		log.Println(err)

		ctx.SetStatusCode(http.StatusInternalServerError)
		ctx.SetBodyString(http.StatusText(http.StatusInternalServerError))

		return
	}

	if err = r.Set(cacheKey, data, time.Duration(config.Cache.RenderCacheDuration)*time.Second); err != nil {
		log.Println(err)

		ctx.SetStatusCode(http.StatusInternalServerError)
		ctx.SetBodyString(http.StatusText(http.StatusInternalServerError))

		return
	}

	if download {
		ctx.Response.Header.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.png"`, user))
	}

	ctx.Response.Header.Set("X-Cache-Hit", "FALSE")
	ctx.SetContentType("image/png")
	ctx.SetBody(data)
}
