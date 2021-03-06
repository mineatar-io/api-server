package routes

import (
	"fmt"
	"log"
	"net/http"

	"github.com/mineatar-io/api-server/src/util"
	"github.com/mineatar-io/skin-render"
	"github.com/valyala/fasthttp"
)

func FullBodyHandler(ctx *fasthttp.RequestCtx) {
	user := ctx.UserValue("user").(string)

	opts := util.ParseQueryParams(ctx, config.Routes.FullBody)

	uuid, ok, err := util.LookupUUID(user)

	if err != nil {
		util.WriteError(ctx, err, http.StatusInternalServerError)

		return
	}

	if !ok && !opts.Fallback {
		util.WriteError(ctx, nil, http.StatusNotFound)

		return
	}

	cacheKey := fmt.Sprintf("result:fullbody-%d-%t-%s", opts.Scale, opts.Overlay, uuid)

	{
		cache, ok, err := r.GetBytes(cacheKey)

		if err != nil {
			util.WriteError(ctx, err, http.StatusInternalServerError)

			return
		}

		if ok {
			if opts.Download {
				ctx.Response.Header.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.png"`, user))
			}

			if util.Debug {
				log.Printf("Retrieved cache for full body render for '%s'\n", uuid)
			}

			ctx.Response.Header.Set("X-Cache-Hit", "TRUE")
			ctx.SetContentType("image/png")
			ctx.SetBody(cache)

			return
		}
	}

	rawSkin, slim, err := util.GetPlayerSkin(uuid)

	if err != nil {
		util.WriteError(ctx, err, http.StatusInternalServerError)

		return
	}

	render := skin.RenderBody(rawSkin, skin.Options{
		Overlay: opts.Overlay,
		Slim:    slim,
		Scale:   opts.Scale,
	})

	if util.Debug {
		log.Printf("Rendered full body image for '%s'\n", uuid)
	}

	data, err := util.EncodePNG(render)

	if err != nil {
		util.WriteError(ctx, err, http.StatusInternalServerError)

		return
	}

	if err = r.Set(cacheKey, data, config.Cache.RenderCacheDuration); err != nil {
		util.WriteError(ctx, err, http.StatusInternalServerError)

		return
	}

	if opts.Download {
		ctx.Response.Header.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.png"`, user))
	}

	ctx.Response.Header.Set("X-Cache-Hit", "FALSE")
	ctx.SetContentType("image/png")
	ctx.SetBody(data)
}

func FrontBodyHandler(ctx *fasthttp.RequestCtx) {
	user := ctx.UserValue("user").(string)

	opts := util.ParseQueryParams(ctx, config.Routes.FrontBody)

	uuid, ok, err := util.LookupUUID(user)

	if err != nil {
		util.WriteError(ctx, err, http.StatusInternalServerError)

		return
	}

	if !ok && !opts.Fallback {
		util.WriteError(ctx, nil, http.StatusNotFound)

		return
	}

	cacheKey := fmt.Sprintf("result:frontbody-%d-%t-%s", opts.Scale, opts.Overlay, uuid)

	{
		cache, ok, err := r.GetBytes(cacheKey)

		if err != nil {
			util.WriteError(ctx, err, http.StatusInternalServerError)

			return
		}

		if ok {
			if opts.Download {
				ctx.Response.Header.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.png"`, user))
			}

			if util.Debug {
				log.Printf("Retrieved cache for front body render for '%s'\n", uuid)
			}

			ctx.Response.Header.Set("X-Cache-Hit", "TRUE")
			ctx.SetContentType("image/png")
			ctx.SetBody(cache)

			return
		}
	}

	rawSkin, slim, err := util.GetPlayerSkin(uuid)

	if err != nil {
		util.WriteError(ctx, err, http.StatusInternalServerError)

		return
	}

	render := skin.RenderFrontBody(rawSkin, skin.Options{
		Overlay: opts.Overlay,
		Slim:    slim,
		Scale:   opts.Scale,
	})

	if util.Debug {
		log.Printf("Rendered front body image for '%s'\n", uuid)
	}

	data, err := util.EncodePNG(render)

	if err != nil {
		util.WriteError(ctx, err, http.StatusInternalServerError)

		return
	}

	if err = r.Set(cacheKey, data, config.Cache.RenderCacheDuration); err != nil {
		util.WriteError(ctx, err, http.StatusInternalServerError)

		return
	}

	if opts.Download {
		ctx.Response.Header.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.png"`, user))
	}

	ctx.Response.Header.Set("X-Cache-Hit", "FALSE")
	ctx.SetContentType("image/png")
	ctx.SetBody(data)
}

func BackBodyHandler(ctx *fasthttp.RequestCtx) {
	user := ctx.UserValue("user").(string)

	opts := util.ParseQueryParams(ctx, config.Routes.BackBody)

	uuid, ok, err := util.LookupUUID(user)

	if err != nil {
		util.WriteError(ctx, err, http.StatusInternalServerError)

		return
	}

	if !ok && !opts.Fallback {
		util.WriteError(ctx, nil, http.StatusNotFound)

		return
	}

	cacheKey := fmt.Sprintf("result:backbody-%d-%t-%s", opts.Scale, opts.Overlay, uuid)

	{
		cache, ok, err := r.GetBytes(cacheKey)

		if err != nil {
			util.WriteError(ctx, err, http.StatusInternalServerError)

			return
		}

		if ok {
			if opts.Download {
				ctx.Response.Header.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.png"`, user))
			}

			if util.Debug {
				log.Printf("Retrieved cache for back body render for '%s'\n", uuid)
			}

			ctx.Response.Header.Set("X-Cache-Hit", "TRUE")
			ctx.SetContentType("image/png")
			ctx.SetBody(cache)

			return
		}
	}

	rawSkin, slim, err := util.GetPlayerSkin(uuid)

	if err != nil {
		util.WriteError(ctx, err, http.StatusInternalServerError)

		return
	}

	render := skin.RenderBackBody(rawSkin, skin.Options{
		Overlay: opts.Overlay,
		Slim:    slim,
		Scale:   opts.Scale,
	})

	if util.Debug {
		log.Printf("Rendered back body image for '%s'\n", uuid)
	}

	data, err := util.EncodePNG(render)

	if err != nil {
		util.WriteError(ctx, err, http.StatusInternalServerError)

		return
	}

	if err = r.Set(cacheKey, data, config.Cache.RenderCacheDuration); err != nil {
		util.WriteError(ctx, err, http.StatusInternalServerError)

		return
	}

	if opts.Download {
		ctx.Response.Header.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.png"`, user))
	}

	ctx.Response.Header.Set("X-Cache-Hit", "FALSE")
	ctx.SetContentType("image/png")
	ctx.SetBody(data)
}

func LeftBodyHandler(ctx *fasthttp.RequestCtx) {
	user := ctx.UserValue("user").(string)

	opts := util.ParseQueryParams(ctx, config.Routes.LeftBody)

	uuid, ok, err := util.LookupUUID(user)

	if err != nil {
		util.WriteError(ctx, err, http.StatusInternalServerError)

		return
	}

	if !ok && !opts.Fallback {
		util.WriteError(ctx, nil, http.StatusNotFound)

		return
	}

	cacheKey := fmt.Sprintf("result:leftbody-%d-%t-%s", opts.Scale, opts.Overlay, uuid)

	{
		cache, ok, err := r.GetBytes(cacheKey)

		if err != nil {
			util.WriteError(ctx, err, http.StatusInternalServerError)

			return
		}

		if ok {
			if opts.Download {
				ctx.Response.Header.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.png"`, user))
			}

			if util.Debug {
				log.Printf("Retrieved cache for left body render for '%s'\n", uuid)
			}

			ctx.Response.Header.Set("X-Cache-Hit", "TRUE")
			ctx.SetContentType("image/png")
			ctx.SetBody(cache)

			return
		}
	}

	rawSkin, slim, err := util.GetPlayerSkin(uuid)

	if err != nil {
		util.WriteError(ctx, err, http.StatusInternalServerError)

		return
	}

	render := skin.RenderLeftBody(rawSkin, skin.Options{
		Overlay: opts.Overlay,
		Slim:    slim,
		Scale:   opts.Scale,
	})

	if util.Debug {
		log.Printf("Rendered left body image for '%s'\n", uuid)
	}

	data, err := util.EncodePNG(render)

	if err != nil {
		util.WriteError(ctx, err, http.StatusInternalServerError)

		return
	}

	if err = r.Set(cacheKey, data, config.Cache.RenderCacheDuration); err != nil {
		util.WriteError(ctx, err, http.StatusInternalServerError)

		return
	}

	if opts.Download {
		ctx.Response.Header.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.png"`, user))
	}

	ctx.Response.Header.Set("X-Cache-Hit", "FALSE")
	ctx.SetContentType("image/png")
	ctx.SetBody(data)
}

func RightBodyHandler(ctx *fasthttp.RequestCtx) {
	user := ctx.UserValue("user").(string)

	opts := util.ParseQueryParams(ctx, config.Routes.RightBody)

	uuid, ok, err := util.LookupUUID(user)

	if err != nil {
		util.WriteError(ctx, err, http.StatusInternalServerError)

		return
	}

	if !ok && !opts.Fallback {
		util.WriteError(ctx, nil, http.StatusNotFound)

		return
	}

	cacheKey := fmt.Sprintf("result:rightbody-%d-%t-%s", opts.Scale, opts.Overlay, uuid)

	{
		cache, ok, err := r.GetBytes(cacheKey)

		if err != nil {
			util.WriteError(ctx, err, http.StatusInternalServerError)

			return
		}

		if ok {
			if opts.Download {
				ctx.Response.Header.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.png"`, user))
			}

			if util.Debug {
				log.Printf("Retrieved cache for right body render for '%s'\n", uuid)
			}

			ctx.Response.Header.Set("X-Cache-Hit", "TRUE")
			ctx.SetContentType("image/png")
			ctx.SetBody(cache)

			return
		}
	}

	rawSkin, slim, err := util.GetPlayerSkin(uuid)

	if err != nil {
		util.WriteError(ctx, err, http.StatusInternalServerError)

		return
	}

	render := skin.RenderRightBody(rawSkin, skin.Options{
		Overlay: opts.Overlay,
		Slim:    slim,
		Scale:   opts.Scale,
	})

	if util.Debug {
		log.Printf("Rendered right body image for '%s'\n", uuid)
	}

	data, err := util.EncodePNG(render)

	if err != nil {
		util.WriteError(ctx, err, http.StatusInternalServerError)

		return
	}

	if err = r.Set(cacheKey, data, config.Cache.RenderCacheDuration); err != nil {
		util.WriteError(ctx, err, http.StatusInternalServerError)

		return
	}

	if opts.Download {
		ctx.Response.Header.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.png"`, user))
	}

	ctx.Response.Header.Set("X-Cache-Hit", "FALSE")
	ctx.SetContentType("image/png")
	ctx.SetBody(data)
}
