package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/mineatar-io/skin-render"
)

func init() {
	app.Get("/ping", PingHandler)
	app.Get("/list", ListHandler)
	app.Get("/skin/:uuid", SkinHandler)
	app.Get("/face/:uuid", FaceHandler)
	app.Get("/head/:uuid", HeadHandler)
	app.Get("/body/full/:uuid", FullBodyHandler)
	app.Get("/body/front/:uuid", FrontBodyHandler)
	app.Get("/body/back/:uuid", BackBodyHandler)
	app.Get("/body/left/:uuid", LeftBodyHandler)
	app.Get("/body/right/:uuid", RightBodyHandler)
	app.Use(NotFoundHandler)
}

// PingHandler is the API handler used for the `/ping` route.
func PingHandler(ctx *fiber.Ctx) error {
	return ctx.SendStatus(http.StatusOK)
}

// ListHandler is the API handler used for the `/list` route.
func ListHandler(ctx *fiber.Ctx) error {
	authKey := ctx.Get("Authorization")

	if len(authKey) < 1 {
		return ctx.SendStatus(http.StatusUnauthorized)
	}

	if authKey != conf.AuthKey {
		return ctx.SendStatus(http.StatusForbidden)
	}

	result := make([]string, 0)

	var (
		cursor uint64 = 0
		keys   []string
		err    error
	)

	for {
		keys, cursor, err = r.Scan(cursor, "unique:*", 25)

		if err != nil {
			return err
		}

		for _, uuid := range keys {
			result = append(result, strings.TrimPrefix(uuid, "unique:"))
		}

		if cursor == 0 {
			break
		}
	}

	return ctx.JSON(result)
}

// FullBodyHandler is the API handler used for the `/body/full/:uuid` route.
func FullBodyHandler(ctx *fiber.Ctx) error {
	opts := ParseQueryParams(ctx, conf.Routes.FullBody)

	uuid, ok := ParseUUID(ExtractUUID(ctx))

	if !ok {
		return SendUsernameDeprecation(ctx)
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

// FrontBodyHandler is the API handler used for the `/body/front/:uuid` route.
func FrontBodyHandler(ctx *fiber.Ctx) error {
	opts := ParseQueryParams(ctx, conf.Routes.FrontBody)

	uuid, ok := ParseUUID(ExtractUUID(ctx))

	if !ok {
		return SendUsernameDeprecation(ctx)
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

// BackBodyHandler is the API handler used for the `/body/back/:uuid` route.
func BackBodyHandler(ctx *fiber.Ctx) error {
	opts := ParseQueryParams(ctx, conf.Routes.BackBody)

	uuid, ok := ParseUUID(ExtractUUID(ctx))

	if !ok {
		return SendUsernameDeprecation(ctx)
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

// LeftBodyHandler is the API handler used for the `/body/left/:uuid` route.
func LeftBodyHandler(ctx *fiber.Ctx) error {
	opts := ParseQueryParams(ctx, conf.Routes.LeftBody)

	uuid, ok := ParseUUID(ExtractUUID(ctx))

	if !ok {
		return SendUsernameDeprecation(ctx)
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

// RightBodyHandler is the API handler used for the `/body/right/:uuid` route.
func RightBodyHandler(ctx *fiber.Ctx) error {
	opts := ParseQueryParams(ctx, conf.Routes.RightBody)

	uuid, ok := ParseUUID(ExtractUUID(ctx))

	if !ok {
		return SendUsernameDeprecation(ctx)
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

// FaceHandler is the API handler used for the `/face/:uuid` route.
func FaceHandler(ctx *fiber.Ctx) error {
	opts := ParseQueryParams(ctx, conf.Routes.Face)

	uuid, ok := ParseUUID(ExtractUUID(ctx))

	if !ok {
		return SendUsernameDeprecation(ctx)
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

// HeadHandler is the API handler used for the `/head/:uuid` route.
func HeadHandler(ctx *fiber.Ctx) error {
	opts := ParseQueryParams(ctx, conf.Routes.Head)

	uuid, ok := ParseUUID(ExtractUUID(ctx))

	if !ok {
		return SendUsernameDeprecation(ctx)
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

// SkinHandler is the API handler used for the `/skin/:uuid` route.
func SkinHandler(ctx *fiber.Ctx) error {
	opts := ParseQueryParams(ctx, conf.Routes.RawSkin)

	uuid, ok := ParseUUID(ctx.Params("uuid"))

	if !ok {
		return SendUsernameDeprecation(ctx)
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

// NotFoundHandler is the API handler used for any requests that do not match an existing route.
func NotFoundHandler(ctx *fiber.Ctx) error {
	return ctx.SendStatus(http.StatusNotFound)
}
