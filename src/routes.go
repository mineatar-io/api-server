package main

import (
	"fmt"
	"net/http"
	"strconv"

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
	return ctx.SendStatus(http.StatusOK)
}

func FullBodyHandler(ctx *fiber.Ctx) error {
	opts := ParseQueryParams(ctx, config.Routes.FullBody)

	uuid, ok, err := LookupUUID(ParseUserParam(ctx))

	if err != nil {
		return err
	}

	if !ok && !opts.Fallback {
		return ctx.SendStatus(http.StatusNotFound)
	}

	cacheKey := fmt.Sprintf("result:fullbody-%d-%t-%s", opts.Scale, opts.Overlay, uuid)

	data, err := r.GetBytes(cacheKey)

	if err != nil {
		return err
	}

	ctx.Set("X-Cache-Hit", strconv.FormatBool(data != nil))

	if data == nil {
		rawSkin, slim, err := GetPlayerSkin(uuid)

		if err != nil {
			return err
		}

		if data, err = EncodePNG(skin.RenderBody(rawSkin, skin.Options{
			Overlay: opts.Overlay,
			Slim:    slim,
			Scale:   opts.Scale,
		})); err != nil {
			return err
		}

		if err = r.Set(cacheKey, data, config.Cache.RenderCacheDuration); err != nil {
			return err
		}
	}

	if opts.Download {
		ctx.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.png"`, uuid))
	}

	return ctx.Type("png").Send(data)
}

func FrontBodyHandler(ctx *fiber.Ctx) error {
	opts := ParseQueryParams(ctx, config.Routes.FrontBody)

	uuid, ok, err := LookupUUID(ParseUserParam(ctx))

	if err != nil {
		return err
	}

	if !ok && !opts.Fallback {
		return ctx.SendStatus(http.StatusNotFound)
	}

	cacheKey := fmt.Sprintf("result:frontbody-%d-%t-%s", opts.Scale, opts.Overlay, uuid)

	data, err := r.GetBytes(cacheKey)

	if err != nil {
		return err
	}

	ctx.Set("X-Cache-Hit", strconv.FormatBool(data != nil))

	if data == nil {
		rawSkin, slim, err := GetPlayerSkin(uuid)

		if err != nil {
			return err
		}

		data, err = EncodePNG(skin.RenderFrontBody(rawSkin, skin.Options{
			Overlay: opts.Overlay,
			Slim:    slim,
			Scale:   opts.Scale,
		}))

		if err != nil {
			return err
		}

		if err = r.Set(cacheKey, data, config.Cache.RenderCacheDuration); err != nil {
			return err
		}
	}

	if opts.Download {
		ctx.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.png"`, uuid))
	}

	return ctx.Type("png").Send(data)
}

func BackBodyHandler(ctx *fiber.Ctx) error {
	opts := ParseQueryParams(ctx, config.Routes.BackBody)

	uuid, ok, err := LookupUUID(ParseUserParam(ctx))

	if err != nil {
		return err
	}

	if !ok && !opts.Fallback {
		return ctx.SendStatus(http.StatusNotFound)
	}

	cacheKey := fmt.Sprintf("result:backbody-%d-%t-%s", opts.Scale, opts.Overlay, uuid)

	data, err := r.GetBytes(cacheKey)

	if err != nil {
		return err
	}

	ctx.Set("X-Cache-Hit", strconv.FormatBool(data != nil))

	if data == nil {
		rawSkin, slim, err := GetPlayerSkin(uuid)

		if err != nil {
			return err
		}

		data, err = EncodePNG(skin.RenderBackBody(rawSkin, skin.Options{
			Overlay: opts.Overlay,
			Slim:    slim,
			Scale:   opts.Scale,
		}))

		if err != nil {
			return err
		}

		if err = r.Set(cacheKey, data, config.Cache.RenderCacheDuration); err != nil {
			return err
		}
	}

	if opts.Download {
		ctx.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.png"`, uuid))
	}

	return ctx.Type("png").Send(data)
}

func LeftBodyHandler(ctx *fiber.Ctx) error {
	opts := ParseQueryParams(ctx, config.Routes.LeftBody)

	uuid, ok, err := LookupUUID(ParseUserParam(ctx))

	if err != nil {
		return err
	}

	if !ok && !opts.Fallback {
		return ctx.SendStatus(http.StatusNotFound)
	}

	cacheKey := fmt.Sprintf("result:leftbody-%d-%t-%s", opts.Scale, opts.Overlay, uuid)

	data, err := r.GetBytes(cacheKey)

	if err != nil {
		return err
	}

	ctx.Set("X-Cache-Hit", strconv.FormatBool(data != nil))

	if data == nil {
		rawSkin, slim, err := GetPlayerSkin(uuid)

		if err != nil {
			return err
		}

		data, err = EncodePNG(skin.RenderLeftBody(rawSkin, skin.Options{
			Overlay: opts.Overlay,
			Slim:    slim,
			Scale:   opts.Scale,
		}))

		if err != nil {
			return err
		}

		if err = r.Set(cacheKey, data, config.Cache.RenderCacheDuration); err != nil {
			return err
		}
	}

	if opts.Download {
		ctx.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.png"`, uuid))
	}

	return ctx.Type("png").Send(data)
}

func RightBodyHandler(ctx *fiber.Ctx) error {
	opts := ParseQueryParams(ctx, config.Routes.RightBody)

	uuid, ok, err := LookupUUID(ParseUserParam(ctx))

	if err != nil {
		return err
	}

	if !ok && !opts.Fallback {
		return ctx.SendStatus(http.StatusNotFound)
	}

	cacheKey := fmt.Sprintf("result:rightbody-%d-%t-%s", opts.Scale, opts.Overlay, uuid)

	data, err := r.GetBytes(cacheKey)

	if err != nil {
		return err
	}

	ctx.Set("X-Cache-Hit", strconv.FormatBool(data != nil))

	if data == nil {
		rawSkin, slim, err := GetPlayerSkin(uuid)

		if err != nil {
			return err
		}

		data, err = EncodePNG(skin.RenderRightBody(rawSkin, skin.Options{
			Overlay: opts.Overlay,
			Slim:    slim,
			Scale:   opts.Scale,
		}))

		if err != nil {
			return err
		}

		if err = r.Set(cacheKey, data, config.Cache.RenderCacheDuration); err != nil {
			return err
		}
	}

	if opts.Download {
		ctx.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.png"`, uuid))
	}

	return ctx.Type("png").Send(data)
}

func FaceHandler(ctx *fiber.Ctx) error {
	opts := ParseQueryParams(ctx, config.Routes.Face)

	uuid, ok, err := LookupUUID(ParseUserParam(ctx))

	if err != nil {
		return err
	}

	if !ok && !opts.Fallback {
		return ctx.SendStatus(http.StatusNotFound)
	}

	cacheKey := fmt.Sprintf("result:face-%d-%t-%s", opts.Scale, opts.Overlay, uuid)

	data, err := r.GetBytes(cacheKey)

	if err != nil {
		return err
	}

	ctx.Set("X-Cache-Hit", strconv.FormatBool(data != nil))

	if data == nil {
		rawSkin, slim, err := GetPlayerSkin(uuid)

		if err != nil {
			return err
		}

		data, err = EncodePNG(skin.RenderFace(rawSkin, skin.Options{
			Overlay: opts.Overlay,
			Slim:    slim,
			Scale:   opts.Scale,
		}))

		if err != nil {
			return err
		}

		if err = r.Set(cacheKey, data, config.Cache.RenderCacheDuration); err != nil {
			return err
		}
	}

	if opts.Download {
		ctx.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.png"`, uuid))
	}

	return ctx.Type("png").Send(data)
}

func HeadHandler(ctx *fiber.Ctx) error {
	opts := ParseQueryParams(ctx, config.Routes.Head)

	uuid, ok, err := LookupUUID(ParseUserParam(ctx))

	if err != nil {
		return err
	}

	if !ok && !opts.Fallback {
		return ctx.SendStatus(http.StatusNotFound)
	}

	cacheKey := fmt.Sprintf("result:head-%d-%t-%s", opts.Scale, opts.Overlay, uuid)

	data, err := r.GetBytes(cacheKey)

	if err != nil {
		return err
	}

	ctx.Set("X-Cache-Hit", strconv.FormatBool(data != nil))

	if data == nil {
		rawSkin, slim, err := GetPlayerSkin(uuid)

		if err != nil {
			return err
		}

		data, err = EncodePNG(skin.RenderHead(rawSkin, skin.Options{
			Overlay: opts.Overlay,
			Slim:    slim,
			Scale:   opts.Scale,
		}))

		if err != nil {
			return err
		}

		if err = r.Set(cacheKey, data, config.Cache.RenderCacheDuration); err != nil {
			return err
		}
	}

	if opts.Download {
		ctx.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.png"`, uuid))
	}

	return ctx.Type("png").Send(data)
}

func SkinHandler(ctx *fiber.Ctx) error {
	opts := ParseQueryParams(ctx, config.Routes.RawSkin)

	uuid, ok, err := LookupUUID(ParseUserParam(ctx))

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
		ctx.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.png"`, uuid))
	}

	return ctx.Type("png").Send(data)
}

func UUIDHandler(ctx *fiber.Ctx) error {
	uuid, ok, err := LookupUUID(ctx.Params("user"))

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
