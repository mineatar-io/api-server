package routes

import (
	"fmt"
	"log"
	"main/src/util"
	"main/src/util/renders"
	"math"
	"net/http"
	"time"

	"github.com/valyala/fasthttp"
)

func FaceHandler(ctx *fasthttp.RequestCtx) {
	user := ctx.UserValue("user").(string)

	download := ctx.QueryArgs().GetBool("download")

	scale, err := ctx.QueryArgs().GetUint("scale")

	if err != nil {
		scale = 4
	}

	scale = int(math.Max(math.Min(float64(scale), MaxScale), MinScale))

	overlay := true

	if ctx.QueryArgs().Has("overlay") {
		overlay = ctx.QueryArgs().GetBool("overlay")
	}

	uuid, err := util.GetUUID(r, user)

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

	skin, slim, err := util.GetPlayerSkin(r, uuid)

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

	if err = r.Set(cacheKey, data, time.Hour*24); err != nil {
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
