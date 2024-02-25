package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func init() {
	app.Use(recover.New())

	app.Use(favicon.New(favicon.Config{
		Data: faviconData,
	}))

	if config.Environment == "development" {
		app.Use(cors.New(cors.Config{
			AllowOrigins:  "*",
			AllowMethods:  "HEAD,OPTIONS,GET",
			ExposeHeaders: "X-Cache-Hit,X-Cache-Time-Remaining",
		}))

		app.Use(logger.New(logger.Config{
			Format:     "${time} ${ip}:${port} -> ${status}: ${method} ${path} (${latency})\n",
			TimeFormat: "2006/01/02 15:04:05",
		}))
	}

	app.Get("/ping", PingHandler)
	app.Get("/skin/:uuid", SkinHandler)
	app.Get("/face/:uuid", FaceHandler)
	app.Get("/head/:uuid", HeadHandler)
	app.Get("/body/full/:uuid", FullBodyHandler)
	app.Get("/body/front/:uuid", FrontBodyHandler)
	app.Get("/body/back/:uuid", BackBodyHandler)
	app.Get("/body/left/:uuid", LeftBodyHandler)
	app.Get("/body/right/:uuid", RightBodyHandler)
}

// PingHandler is the API handler used for the `/ping` route.
func PingHandler(ctx *fiber.Ctx) error {
	return ctx.SendStatus(http.StatusOK)
}

// SkinHandler is the API handler used for the `/skin/:uuid` route.
func SkinHandler(ctx *fiber.Ctx) error {
	opts := ParseQueryParams(ctx, config.Routes.RawSkin)

	if opts == nil {
		return nil
	}

	uuid, ok := ParseUUID(ctx.Params("uuid"))

	if !ok {
		return ctx.Status(http.StatusBadRequest).SendString("Invalid UUID")
	}

	rawSkin, _, err := GetPlayerSkin(uuid)

	if err != nil {
		return err
	}

	data, err := EncodeImage(rawSkin, opts)

	if err != nil {
		return err
	}

	if opts.Download {
		ctx.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.%s"`, uuid, opts.Format))
	}

	return ctx.Type(opts.Format).Send(data)
}

// FaceHandler is the API handler used for the `/face/:uuid` route.
func FaceHandler(ctx *fiber.Ctx) error {
	opts := ParseQueryParams(ctx, config.Routes.Face)

	if opts == nil {
		return nil
	}

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

	return ctx.Type(opts.Format).Send(result)
}

// HeadHandler is the API handler used for the `/head/:uuid` route.
func HeadHandler(ctx *fiber.Ctx) error {
	opts := ParseQueryParams(ctx, config.Routes.Head)

	if opts == nil {
		return nil
	}

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

	return ctx.Type(opts.Format).Send(result)
}

// FullBodyHandler is the API handler used for the `/body/full/:uuid` route.
func FullBodyHandler(ctx *fiber.Ctx) error {
	opts := ParseQueryParams(ctx, config.Routes.FullBody)

	if opts == nil {
		return nil
	}

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

	return ctx.Type(opts.Format).Send(result)
}

// FrontBodyHandler is the API handler used for the `/body/front/:uuid` route.
func FrontBodyHandler(ctx *fiber.Ctx) error {
	opts := ParseQueryParams(ctx, config.Routes.FrontBody)

	if opts == nil {
		return nil
	}

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

	return ctx.Type(opts.Format).Send(result)
}

// BackBodyHandler is the API handler used for the `/body/back/:uuid` route.
func BackBodyHandler(ctx *fiber.Ctx) error {
	opts := ParseQueryParams(ctx, config.Routes.BackBody)

	if opts == nil {
		return nil
	}

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

	return ctx.Type(opts.Format).Send(result)
}

// LeftBodyHandler is the API handler used for the `/body/left/:uuid` route.
func LeftBodyHandler(ctx *fiber.Ctx) error {
	opts := ParseQueryParams(ctx, config.Routes.LeftBody)

	if opts == nil {
		return nil
	}

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

	return ctx.Type(opts.Format).Send(result)
}

// RightBodyHandler is the API handler used for the `/body/right/:uuid` route.
func RightBodyHandler(ctx *fiber.Ctx) error {
	opts := ParseQueryParams(ctx, config.Routes.RightBody)

	if opts == nil {
		return nil
	}

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

	return ctx.Type(opts.Format).Send(result)
}
