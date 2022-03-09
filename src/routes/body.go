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

func FullBodyHandler(ctx *fasthttp.RequestCtx) {
	user := ctx.UserValue("user").(string)

	download := ctx.QueryArgs().GetBool("download")

	scale, err := ctx.QueryArgs().GetUint("scale")

	if err != nil {
		scale = 4
	}

	scale = int(math.Max(math.Min(float64(scale), MaxScaleFullBody), MinScale))

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

	cacheKey := fmt.Sprintf("result:fullbody-%d-%t-%s", scale, overlay, uuid)

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

	render := renders.RenderBody(skin, renders.RenderOptions{
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

func FrontBodyHandler(ctx *fasthttp.RequestCtx) {
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

	cacheKey := fmt.Sprintf("result:frontbody-%d-%t-%s", scale, overlay, uuid)

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

	render := renders.RenderFrontBody(skin, renders.RenderOptions{
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

func BackBodyHandler(ctx *fasthttp.RequestCtx) {
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

	cacheKey := fmt.Sprintf("result:backbody-%d-%t-%s", scale, overlay, uuid)

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

	render := renders.RenderBackBody(skin, renders.RenderOptions{
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

func LeftBodyHandler(ctx *fasthttp.RequestCtx) {
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

	cacheKey := fmt.Sprintf("result:leftbody-%d-%t-%s", scale, overlay, uuid)

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

	render := renders.RenderLeftBody(skin, renders.RenderOptions{
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

func RightBodyHandler(ctx *fasthttp.RequestCtx) {
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

	cacheKey := fmt.Sprintf("result:rightbody-%d-%t-%s", scale, overlay, uuid)

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

	render := renders.RenderRightBody(skin, renders.RenderOptions{
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
