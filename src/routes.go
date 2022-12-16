package main

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/mineatar-io/skin-render"
)

func init() {
	app.Get("/ping", PingHandler)
	app.Get("/uuid/:user", UUIDHandler)
	app.Get("/skin/:user", SkinHandler)
	app.Get("/face/:user", FaceHandler)
	app.Get("/head/:user", HeadHandler)
	app.Get("/body/full/:user", FullBodyHandler)
	app.Get("/body/front/:user", FrontBodyHandler)
	app.Get("/body/back/:user", BackBodyHandler)
	app.Get("/body/left/:user", LeftBodyHandler)
	app.Get("/body/right/:user", RightBodyHandler)
	app.Use(NotFoundHandler)
}

func PingHandler(ctx *fiber.Ctx) error {
	return ctx.SendString("Pong!")
}

func FullBodyHandler(ctx *fiber.Ctx) error {
	user := ctx.Params("user")

	opts := ParseQueryParams(ctx, config.Routes.FullBody)

	uuid, ok, err := LookupUUID(user)

	if err != nil {
		return err
	}

	if !ok && !opts.Fallback {
		return ctx.SendStatus(http.StatusNotFound)
	}

	cacheKey := fmt.Sprintf("result:fullbody-%d-%t-%s", opts.Scale, opts.Overlay, uuid)

	{
		cache, ok, err := r.GetBytes(cacheKey)

		if err != nil {
			return err
		}

		if ok {
			if opts.Download {
				ctx.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.png"`, user))
			}

			ctx.Set("X-Cache-Hit", "TRUE")

			return ctx.Type("png").Send(cache)
		}
	}

	rawSkin, slim, err := GetPlayerSkin(uuid)

	if err != nil {
		return err
	}

	render := skin.RenderBody(rawSkin, skin.Options{
		Overlay: opts.Overlay,
		Slim:    slim,
		Scale:   opts.Scale,
	})

	data, err := EncodePNG(render)

	if err != nil {
		return err
	}

	if err = r.Set(cacheKey, data, config.Cache.RenderCacheDuration); err != nil {
		return err
	}

	if opts.Download {
		ctx.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.png"`, user))
	}

	ctx.Set("X-Cache-Hit", "FALSE")

	return ctx.Type("png").Send(data)
}

func FrontBodyHandler(ctx *fiber.Ctx) error {
	user := ctx.Params("user")

	opts := ParseQueryParams(ctx, config.Routes.FrontBody)

	uuid, ok, err := LookupUUID(user)

	if err != nil {
		return err
	}

	if !ok && !opts.Fallback {
		return ctx.SendStatus(http.StatusNotFound)
	}

	cacheKey := fmt.Sprintf("result:frontbody-%d-%t-%s", opts.Scale, opts.Overlay, uuid)

	{
		cache, ok, err := r.GetBytes(cacheKey)

		if err != nil {
			return err
		}

		if ok {
			if opts.Download {
				ctx.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.png"`, user))
			}

			ctx.Set("X-Cache-Hit", "TRUE")

			return ctx.Type("png").Send(cache)
		}
	}

	rawSkin, slim, err := GetPlayerSkin(uuid)

	if err != nil {
		return err
	}

	render := skin.RenderFrontBody(rawSkin, skin.Options{
		Overlay: opts.Overlay,
		Slim:    slim,
		Scale:   opts.Scale,
	})

	data, err := EncodePNG(render)

	if err != nil {
		return err
	}

	if err = r.Set(cacheKey, data, config.Cache.RenderCacheDuration); err != nil {
		return err
	}

	if opts.Download {
		ctx.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.png"`, user))
	}

	ctx.Set("X-Cache-Hit", "FALSE")

	return ctx.Type("png").Send(data)
}

func BackBodyHandler(ctx *fiber.Ctx) error {
	user := ctx.Params("user")

	opts := ParseQueryParams(ctx, config.Routes.BackBody)

	uuid, ok, err := LookupUUID(user)

	if err != nil {
		return err
	}

	if !ok && !opts.Fallback {
		return ctx.SendStatus(http.StatusNotFound)
	}

	cacheKey := fmt.Sprintf("result:backbody-%d-%t-%s", opts.Scale, opts.Overlay, uuid)

	{
		cache, ok, err := r.GetBytes(cacheKey)

		if err != nil {
			return err
		}

		if ok {
			if opts.Download {
				ctx.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.png"`, user))
			}

			ctx.Set("X-Cache-Hit", "TRUE")

			return ctx.Type("png").Send(cache)
		}
	}

	rawSkin, slim, err := GetPlayerSkin(uuid)

	if err != nil {
		return err
	}

	render := skin.RenderBackBody(rawSkin, skin.Options{
		Overlay: opts.Overlay,
		Slim:    slim,
		Scale:   opts.Scale,
	})

	data, err := EncodePNG(render)

	if err != nil {
		return err
	}

	if err = r.Set(cacheKey, data, config.Cache.RenderCacheDuration); err != nil {
		return err
	}

	if opts.Download {
		ctx.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.png"`, user))
	}

	ctx.Set("X-Cache-Hit", "FALSE")

	return ctx.Type("png").Send(data)
}

func LeftBodyHandler(ctx *fiber.Ctx) error {
	user := ctx.Params("user")

	opts := ParseQueryParams(ctx, config.Routes.LeftBody)

	uuid, ok, err := LookupUUID(user)

	if err != nil {
		return err
	}

	if !ok && !opts.Fallback {
		return ctx.SendStatus(http.StatusNotFound)
	}

	cacheKey := fmt.Sprintf("result:leftbody-%d-%t-%s", opts.Scale, opts.Overlay, uuid)

	{
		cache, ok, err := r.GetBytes(cacheKey)

		if err != nil {
			return err
		}

		if ok {
			if opts.Download {
				ctx.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.png"`, user))
			}

			ctx.Set("X-Cache-Hit", "TRUE")

			return ctx.Type("png").Send(cache)
		}
	}

	rawSkin, slim, err := GetPlayerSkin(uuid)

	if err != nil {
		return err
	}

	render := skin.RenderLeftBody(rawSkin, skin.Options{
		Overlay: opts.Overlay,
		Slim:    slim,
		Scale:   opts.Scale,
	})

	data, err := EncodePNG(render)

	if err != nil {
		return err
	}

	if err = r.Set(cacheKey, data, config.Cache.RenderCacheDuration); err != nil {
		return err
	}

	if opts.Download {
		ctx.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.png"`, user))
	}

	ctx.Set("X-Cache-Hit", "FALSE")

	return ctx.Type("png").Send(data)
}

func RightBodyHandler(ctx *fiber.Ctx) error {
	user := ctx.Params("user")

	opts := ParseQueryParams(ctx, config.Routes.RightBody)

	uuid, ok, err := LookupUUID(user)

	if err != nil {
		return err
	}

	if !ok && !opts.Fallback {
		return ctx.SendStatus(http.StatusNotFound)
	}

	cacheKey := fmt.Sprintf("result:rightbody-%d-%t-%s", opts.Scale, opts.Overlay, uuid)

	{
		cache, ok, err := r.GetBytes(cacheKey)

		if err != nil {
			return err
		}

		if ok {
			if opts.Download {
				ctx.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.png"`, user))
			}

			ctx.Set("X-Cache-Hit", "TRUE")

			return ctx.Type("png").Send(cache)
		}
	}

	rawSkin, slim, err := GetPlayerSkin(uuid)

	if err != nil {
		return err
	}

	render := skin.RenderRightBody(rawSkin, skin.Options{
		Overlay: opts.Overlay,
		Slim:    slim,
		Scale:   opts.Scale,
	})

	data, err := EncodePNG(render)

	if err != nil {
		return err
	}

	if err = r.Set(cacheKey, data, config.Cache.RenderCacheDuration); err != nil {
		return err
	}

	if opts.Download {
		ctx.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.png"`, user))
	}

	ctx.Set("X-Cache-Hit", "FALSE")

	return ctx.Type("png").Send(data)
}

func FaceHandler(ctx *fiber.Ctx) error {
	user := ctx.Params("user")

	opts := ParseQueryParams(ctx, config.Routes.Face)

	uuid, ok, err := LookupUUID(user)

	if err != nil {
		return err
	}

	if !ok && !opts.Fallback {
		return ctx.SendStatus(http.StatusNotFound)
	}

	cacheKey := fmt.Sprintf("result:face-%d-%t-%s", opts.Scale, opts.Overlay, uuid)

	{
		cache, ok, err := r.GetBytes(cacheKey)

		if err != nil {
			return err
		}

		if ok {
			if opts.Download {
				ctx.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.png"`, user))
			}

			ctx.Set("X-Cache-Hit", "TRUE")

			return ctx.Type("png").Send(cache)
		}
	}

	rawSkin, slim, err := GetPlayerSkin(uuid)

	if err != nil {
		return err
	}

	render := skin.RenderFace(rawSkin, skin.Options{
		Overlay: opts.Overlay,
		Slim:    slim,
		Scale:   opts.Scale,
	})

	data, err := EncodePNG(render)

	if err != nil {
		return err
	}

	if err = r.Set(cacheKey, data, config.Cache.RenderCacheDuration); err != nil {
		return err
	}

	if opts.Download {
		ctx.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.png"`, user))
	}

	ctx.Set("X-Cache-Hit", "FALSE")

	return ctx.Type("png").Send(data)
}

func HeadHandler(ctx *fiber.Ctx) error {
	user := ctx.Params("user")

	opts := ParseQueryParams(ctx, config.Routes.Head)

	uuid, ok, err := LookupUUID(user)

	if err != nil {
		return err
	}

	if !ok && !opts.Fallback {
		return ctx.SendStatus(http.StatusNotFound)
	}

	cacheKey := fmt.Sprintf("result:head-%d-%t-%s", opts.Scale, opts.Overlay, uuid)

	{
		cache, ok, err := r.GetBytes(cacheKey)

		if err != nil {
			return err
		}

		if ok {
			if opts.Download {
				ctx.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.png"`, user))
			}

			ctx.Set("X-Cache-Hit", "TRUE")

			return ctx.Type("png").Send(cache)
		}
	}

	rawSkin, slim, err := GetPlayerSkin(uuid)

	if err != nil {
		return err
	}

	render := skin.RenderHead(rawSkin, skin.Options{
		Overlay: opts.Overlay,
		Slim:    slim,
		Scale:   opts.Scale,
	})

	data, err := EncodePNG(render)

	if err != nil {
		return err
	}

	if err = r.Set(cacheKey, data, config.Cache.RenderCacheDuration); err != nil {
		return err
	}

	if opts.Download {
		ctx.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.png"`, user))
	}

	ctx.Set("X-Cache-Hit", "FALSE")

	return ctx.Type("png").Send(data)
}

func SkinHandler(ctx *fiber.Ctx) error {
	user := ctx.Params("user")

	opts := ParseQueryParams(ctx, config.Routes.RawSkin)

	uuid, ok, err := LookupUUID(user)

	if err != nil {
		return err
	}

	if !ok && !opts.Fallback {
		return ctx.SendStatus(http.StatusNotFound)
	}

	rawSkin, _, err := GetPlayerSkin(uuid)

	if err != nil {
		return err
	}

	data, err := EncodePNG(rawSkin)

	if err != nil {
		return err
	}

	if opts.Download {
		ctx.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.png"`, user))
	}

	return ctx.Type("png").Send(data)
}

func UUIDHandler(ctx *fiber.Ctx) error {
	user := ctx.Params("user")

	uuid, ok, err := LookupUUID(user)

	if err != nil {
		return err
	}

	if !ok {
		return ctx.SendStatus(http.StatusNotFound)
	}

	return ctx.SendString(uuid)
}

func NotFoundHandler(ctx *fiber.Ctx) error {
	return ctx.SendStatus(http.StatusNotFound)
}
