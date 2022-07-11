package routes

import (
	"fmt"
	"log"
	"net/http"

	"github.com/mineatar-io/api-server/src/util"
	"github.com/mineatar-io/skin-render"
	"github.com/valyala/fasthttp"
)

func HeadHandler(ctx *fasthttp.RequestCtx) {
	user := ctx.UserValue("user").(string)

	opts := util.ParseQueryParams(ctx, config.Routes.Head)

	uuid, ok, err := util.LookupUUID(user)

	if err != nil {
		util.WriteError(ctx, err, http.StatusInternalServerError)

		return
	}

	if !ok && !opts.Fallback {
		util.WriteError(ctx, nil, http.StatusNotFound)

		return
	}

	cacheKey := fmt.Sprintf("result:head-%d-%t-%s", opts.Scale, opts.Overlay, uuid)

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
				log.Printf("Retrieved cache for head render for '%s'\n", uuid)
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

	render := skin.RenderHead(rawSkin, skin.Options{
		Overlay: opts.Overlay,
		Slim:    slim,
		Scale:   opts.Scale,
	})

	if util.Debug {
		log.Printf("Rendered head image for '%s'\n", uuid)
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
