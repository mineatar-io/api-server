package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func init() {
	app.Get("/ping", PingHandler)
	app.Get("/favicon.ico", FaviconHandler)
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

// FaviconHandler serves the favicon.ico file to any users that visit the API using a browser.
func FaviconHandler(ctx *fiber.Ctx) error {
	return ctx.Type("ico").Send(favicon)
}

// ListHandler is the API handler used for the `/list` route.
func ListHandler(ctx *fiber.Ctx) error {
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

// SkinHandler is the API handler used for the `/skin/:uuid` route.
func SkinHandler(ctx *fiber.Ctx) error {
	opts := ParseQueryParams(ctx, config.Routes.RawSkin)

	uuid, ok := ParseUUID(ctx.Params("uuid"))

	if !ok {
		return ctx.Status(http.StatusBadRequest).SendString("Invalid UUID")
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

// FaceHandler is the API handler used for the `/face/:uuid` route.
func FaceHandler(ctx *fiber.Ctx) error {
	opts := ParseQueryParams(ctx, config.Routes.Face)

	uuid, ok := ParseUUID(ExtractUUID(ctx))

	if !ok {
		return ctx.Status(http.StatusBadRequest).SendString("Invalid UUID")
	}

	rawSkin, isSlim, err := GetPlayerSkin(uuid)

	if err != nil {
		return err
	}

	result, cache, err := Render(RenderTypeFace, uuid, rawSkin, isSlim, opts)

	if err != nil {
		return err
	}

	ctx.Set("X-Cache-Hit", strconv.FormatBool(cache))

	if opts.Download {
		ctx.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.png"`, uuid))
	}

	return ctx.Type("png").Send(result)
}

// HeadHandler is the API handler used for the `/head/:uuid` route.
func HeadHandler(ctx *fiber.Ctx) error {
	opts := ParseQueryParams(ctx, config.Routes.Head)

	uuid, ok := ParseUUID(ExtractUUID(ctx))

	if !ok {
		return ctx.Status(http.StatusBadRequest).SendString("Invalid UUID")
	}

	rawSkin, isSlim, err := GetPlayerSkin(uuid)

	if err != nil {
		return err
	}

	result, cache, err := Render(RenderTypeHead, uuid, rawSkin, isSlim, opts)

	if err != nil {
		return err
	}

	ctx.Set("X-Cache-Hit", strconv.FormatBool(cache))

	if opts.Download {
		ctx.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.png"`, uuid))
	}

	return ctx.Type("png").Send(result)
}

// FullBodyHandler is the API handler used for the `/body/full/:uuid` route.
func FullBodyHandler(ctx *fiber.Ctx) error {
	opts := ParseQueryParams(ctx, config.Routes.FullBody)

	uuid, ok := ParseUUID(ExtractUUID(ctx))

	if !ok {
		return ctx.Status(http.StatusBadRequest).SendString("Invalid UUID")
	}

	rawSkin, isSlim, err := GetPlayerSkin(uuid)

	if err != nil {
		return err
	}

	result, cache, err := Render(RenderTypeFullBody, uuid, rawSkin, isSlim, opts)

	if err != nil {
		return err
	}

	ctx.Set("X-Cache-Hit", strconv.FormatBool(cache))

	if opts.Download {
		ctx.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.png"`, uuid))
	}

	return ctx.Type("png").Send(result)
}

// FrontBodyHandler is the API handler used for the `/body/front/:uuid` route.
func FrontBodyHandler(ctx *fiber.Ctx) error {
	opts := ParseQueryParams(ctx, config.Routes.FrontBody)

	uuid, ok := ParseUUID(ExtractUUID(ctx))

	if !ok {
		return ctx.Status(http.StatusBadRequest).SendString("Invalid UUID")
	}

	rawSkin, isSlim, err := GetPlayerSkin(uuid)

	if err != nil {
		return err
	}

	result, cache, err := Render(RenderTypeFrontBody, uuid, rawSkin, isSlim, opts)

	if err != nil {
		return err
	}

	ctx.Set("X-Cache-Hit", strconv.FormatBool(cache))

	if opts.Download {
		ctx.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.png"`, uuid))
	}

	return ctx.Type("png").Send(result)
}

// BackBodyHandler is the API handler used for the `/body/back/:uuid` route.
func BackBodyHandler(ctx *fiber.Ctx) error {
	opts := ParseQueryParams(ctx, config.Routes.BackBody)

	uuid, ok := ParseUUID(ExtractUUID(ctx))

	if !ok {
		return ctx.Status(http.StatusBadRequest).SendString("Invalid UUID")
	}

	rawSkin, isSlim, err := GetPlayerSkin(uuid)

	if err != nil {
		return err
	}

	result, cache, err := Render(RenderTypeBackBody, uuid, rawSkin, isSlim, opts)

	if err != nil {
		return err
	}

	ctx.Set("X-Cache-Hit", strconv.FormatBool(cache))

	if opts.Download {
		ctx.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.png"`, uuid))
	}

	return ctx.Type("png").Send(result)
}

// LeftBodyHandler is the API handler used for the `/body/left/:uuid` route.
func LeftBodyHandler(ctx *fiber.Ctx) error {
	opts := ParseQueryParams(ctx, config.Routes.LeftBody)

	uuid, ok := ParseUUID(ExtractUUID(ctx))

	if !ok {
		return ctx.Status(http.StatusBadRequest).SendString("Invalid UUID")
	}

	rawSkin, isSlim, err := GetPlayerSkin(uuid)

	if err != nil {
		return err
	}

	result, cache, err := Render(RenderTypeLeftBody, uuid, rawSkin, isSlim, opts)

	if err != nil {
		return err
	}

	ctx.Set("X-Cache-Hit", strconv.FormatBool(cache))

	if opts.Download {
		ctx.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.png"`, uuid))
	}

	return ctx.Type("png").Send(result)
}

// RightBodyHandler is the API handler used for the `/body/right/:uuid` route.
func RightBodyHandler(ctx *fiber.Ctx) error {
	opts := ParseQueryParams(ctx, config.Routes.RightBody)

	uuid, ok := ParseUUID(ExtractUUID(ctx))

	if !ok {
		return ctx.Status(http.StatusBadRequest).SendString("Invalid UUID")
	}

	rawSkin, isSlim, err := GetPlayerSkin(uuid)

	if err != nil {
		return err
	}

	result, cache, err := Render(RenderTypeRightBody, uuid, rawSkin, isSlim, opts)

	if err != nil {
		return err
	}

	ctx.Set("X-Cache-Hit", strconv.FormatBool(cache))

	if opts.Download {
		ctx.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.png"`, uuid))
	}

	return ctx.Type("png").Send(result)
}

// NotFoundHandler is the API handler used for any requests that do not match an existing route.
func NotFoundHandler(ctx *fiber.Ctx) error {
	return ctx.SendStatus(http.StatusNotFound)
}
