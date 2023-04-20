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

// PingHandler is the API handler used for the `/ping` route.
func PingHandler(ctx *fiber.Ctx) error {
	return ctx.SendStatus(http.StatusOK)
}

// FullBodyHandler is the API handler used for the `/body/full/:user` route.
func FullBodyHandler(ctx *fiber.Ctx) error {
	opts := ParseQueryParams(ctx, conf.Routes.FullBody)

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

		if err = r.Set(cacheKey, data, conf.Cache.RenderCacheDuration); err != nil {
			return err
		}
	}

	if opts.Download {
		ctx.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.png"`, uuid))
	}

	return ctx.Type("png").Send(data)
}

// FrontBodyHandler is the API handler used for the `/body/front/:user` route.
func FrontBodyHandler(ctx *fiber.Ctx) error {
	opts := ParseQueryParams(ctx, conf.Routes.FrontBody)

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

		if err = r.Set(cacheKey, data, conf.Cache.RenderCacheDuration); err != nil {
			return err
		}
	}

	if opts.Download {
		ctx.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.png"`, uuid))
	}

	return ctx.Type("png").Send(data)
}

// BackBodyHandler is the API handler used for the `/body/back/:user` route.
func BackBodyHandler(ctx *fiber.Ctx) error {
	opts := ParseQueryParams(ctx, conf.Routes.BackBody)

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

		if err = r.Set(cacheKey, data, conf.Cache.RenderCacheDuration); err != nil {
			return err
		}
	}

	if opts.Download {
		ctx.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.png"`, uuid))
	}

	return ctx.Type("png").Send(data)
}

// LeftBodyHandler is the API handler used for the `/body/left/:user` route.
func LeftBodyHandler(ctx *fiber.Ctx) error {
	opts := ParseQueryParams(ctx, conf.Routes.LeftBody)

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

		if err = r.Set(cacheKey, data, conf.Cache.RenderCacheDuration); err != nil {
			return err
		}
	}

	if opts.Download {
		ctx.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.png"`, uuid))
	}

	return ctx.Type("png").Send(data)
}

// RightBodyHandler is the API handler used for the `/body/right/:user` route.
func RightBodyHandler(ctx *fiber.Ctx) error {
	opts := ParseQueryParams(ctx, conf.Routes.RightBody)

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

		if err = r.Set(cacheKey, data, conf.Cache.RenderCacheDuration); err != nil {
			return err
		}
	}

	if opts.Download {
		ctx.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.png"`, uuid))
	}

	return ctx.Type("png").Send(data)
}

// FaceHandler is the API handler used for the `/face/:user` route.
func FaceHandler(ctx *fiber.Ctx) error {
	opts := ParseQueryParams(ctx, conf.Routes.Face)

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

		if err = r.Set(cacheKey, data, conf.Cache.RenderCacheDuration); err != nil {
			return err
		}
	}

	if opts.Download {
		ctx.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.png"`, uuid))
	}

	return ctx.Type("png").Send(data)
}

// HeadHandler is the API handler used for the `/head/:user` route.
func HeadHandler(ctx *fiber.Ctx) error {
	opts := ParseQueryParams(ctx, conf.Routes.Head)

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

		if err = r.Set(cacheKey, data, conf.Cache.RenderCacheDuration); err != nil {
			return err
		}
	}

	if opts.Download {
		ctx.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.png"`, uuid))
	}

	return ctx.Type("png").Send(data)
}

// SkinHandler is the API handler used for the `/skin/:user` route.
func SkinHandler(ctx *fiber.Ctx) error {
	opts := ParseQueryParams(ctx, conf.Routes.RawSkin)

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

// UUIDHandler is the API handler used for the `/uuid/:user` route.
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

// NotFoundHandler is the API handler used for any requests that do not match an existing route.
func NotFoundHandler(ctx *fiber.Ctx) error {
	return ctx.SendStatus(http.StatusNotFound)
}
