package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/mineatar-io/skin-render"
)

var (
	//go:embed favicon.ico
	favicon []byte

	RenderTypeFullBody  = "fullbody"
	RenderTypeFrontBody = "frontbody"
	RenderTypeBackBody  = "backbody"
	RenderTypeLeftBody  = "leftbody"
	RenderTypeRightBody = "rightbody"
	RenderTypeFace      = "face"
	RenderTypeHead      = "head"
)

// QueryParams is used by most all API routes as options for how the image should be rendered, or how errors should be handled.
type QueryParams struct {
	Scale    int  `query:"scale"`
	Download bool `query:"download"`
	Overlay  bool `query:"overlay"`
}

// PointerOf returns the value of the first argument as a pointer.
func PointerOf[T any](v T) *T {
	return &v
}

// Clamp clamps the input value between the minimum and maximum values.
// This method is preferred over `math.Min()` and `math.Max()` to prevent any type coercion between floats and integers.
func Clamp[T int | uint | int8 | uint8 | int16 | uint16 | int32 | uint32 | int64 | uint64](value, min, max T) T {
	if value > max {
		return max
	}

	if value < min {
		return min
	}

	return value
}

// Render will render the image using the specified details and return the result.
func Render(renderType, uuid string, rawSkin *image.NRGBA, isSlim bool, opts *QueryParams) ([]byte, bool, error) {
	if conf.Cache.EnableLocks {
		mutex := r.NewMutex(fmt.Sprintf("render-lock:%s-%d-%t-%s", renderType, opts.Scale, opts.Overlay, uuid))
		mutex.Lock()

		defer mutex.Unlock()
	}

	cache, err := GetCachedRenderResult(renderType, uuid, opts)

	if err != nil {
		return nil, false, err
	}

	if cache != nil {
		if conf.Environment == "development" {
			log.Printf("Retrieved render from cache (type=%s, uuid=%s, slim=%v, scale=%d)\n", renderType, uuid, isSlim, opts.Scale)
		}

		return cache, true, nil
	}

	var (
		result     *image.NRGBA
		renderOpts skin.Options = skin.Options{
			Overlay: opts.Overlay,
			Slim:    isSlim,
			Scale:   opts.Scale,
		}
	)

	switch renderType {
	case RenderTypeFullBody:
		{
			result = skin.RenderBody(rawSkin, renderOpts)

			break
		}
	case RenderTypeFrontBody:
		{
			result = skin.RenderFrontBody(rawSkin, renderOpts)

			break
		}
	case RenderTypeBackBody:
		{
			result = skin.RenderBackBody(rawSkin, renderOpts)

			break
		}
	case RenderTypeLeftBody:
		{
			result = skin.RenderLeftBody(rawSkin, renderOpts)

			break
		}
	case RenderTypeRightBody:
		{
			result = skin.RenderRightBody(rawSkin, renderOpts)

			break
		}
	case RenderTypeHead:
		{
			result = skin.RenderHead(rawSkin, renderOpts)

			break
		}
	case RenderTypeFace:
		{
			result = skin.RenderFace(rawSkin, renderOpts)

			break
		}
	default:
		panic(fmt.Errorf("unknown render type: %s", renderType))
	}

	data, err := EncodePNG(result)

	if err != nil {
		return nil, false, err
	}

	if err = SetCachedRenderResult(renderType, uuid, opts, data); err != nil {
		return nil, false, err
	}

	if conf.Environment == "development" {
		log.Printf("Rendered image (type=%s, uuid=%s, slim=%v, scale=%d)\n", renderType, uuid, isSlim, opts.Scale)
	}

	return data, false, nil
}

// ExtractUUID returns the UUID from the route param, allowing values such as "<uuid>.png" to be returned as "<uuid>".
func ExtractUUID(ctx *fiber.Ctx) string {
	return strings.Split(ctx.Params("uuid"), ".")[0]
}

// ParseUUID parses the UUID given by the route parameters, and returns a boolean if the UUID is valid.
func ParseUUID(value string) (string, bool) {
	value = strings.ToLower(strings.ReplaceAll(value, "-", ""))

	if len(value) != 32 {
		return "", false
	}

	return value, true
}

// FetchImage fetches the image by the URL and returns it as a parsed image.
func FetchImage(url string) (*image.NRGBA, error) {
	resp, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	img, format, err := image.Decode(resp.Body)

	if err != nil {
		return nil, err
	}

	if format != "NRGBA" {
		outputImg := image.NewNRGBA(img.Bounds())

		draw.Draw(outputImg, img.Bounds(), img, image.Pt(0, 0), draw.Src)

		return outputImg, nil
	}

	return img.(*image.NRGBA), nil
}

// GetPlayerSkin fetches the skin of the Minecraft player by the UUID.
func GetPlayerSkin(uuid string) (*image.NRGBA, bool, error) {
	if conf.Cache.EnableLocks {
		mutex := r.NewMutex(fmt.Sprintf("skin-lock:%s", uuid))
		mutex.Lock()

		defer mutex.Unlock()
	}

	// Get skin from cache, and return if it exists
	if conf.Cache.SkinCacheDuration != nil {
		rawSkin, slim, err := GetCachedSkin(uuid)

		if err != nil {
			return nil, false, err
		}

		if rawSkin != nil {
			return rawSkin, slim, nil
		}
	}

	var (
		err              error                     = nil
		skinImage        *image.NRGBA              = nil
		rawSkin          []byte                    = nil
		isSlim           bool                      = skin.IsSlimFromUUID(uuid)
		profile          *MinecraftProfile         = nil
		rawTextures      string                    = ""
		texturesProperty *MinecraftDecodedTextures = nil
	)

	// Get the textures metadata from Mojang about the Minecraft player
	{
		if profile, err = GetMinecraftProfile(uuid); err != nil {
			return skin.GetDefaultSkin(isSlim), true, nil
		}

		if profile == nil {
			return skin.GetDefaultSkin(isSlim), isSlim, nil
		}

		if err = r.Set(fmt.Sprintf("unique:%s", profile.UUID), "0", 0); err != nil {
			return nil, false, err
		}
	}

	// Locate the skin information within the Minecraft profile properties
	{
		for _, property := range profile.Properties {
			if property.Name != "textures" {
				continue
			}

			rawTextures = property.Value
		}

		if len(rawTextures) < 1 {
			return skin.GetDefaultSkin(isSlim), isSlim, nil
		}
	}

	// Decode the raw textures value returned from the player's properties
	{
		if texturesProperty, err = GetDecodedTexturesValue(rawTextures); err != nil {
			return nil, false, err
		}

		if len(texturesProperty.Textures.Skin.URL) < 1 {
			return skin.GetDefaultSkin(isSlim), isSlim, nil
		}

		isSlim = texturesProperty.Textures.Skin.Metadata.Model == "slim"
	}

	// Fetch the raw skin image from the Mojang API
	{
		if skinImage, err = FetchImage(texturesProperty.Textures.Skin.URL); err != nil {
			return nil, false, err
		}

		if rawSkin, err = EncodePNG(skinImage); err != nil {
			return nil, false, err
		}
	}

	// Put the skin into cache so it can be used for future requests
	if conf.Cache.SkinCacheDuration != nil {
		if err = r.Set(fmt.Sprintf("skin:%s", uuid), rawSkin, *conf.Cache.SkinCacheDuration); err != nil {
			return nil, false, err
		}

		if isSlim {
			if err = r.Set(fmt.Sprintf("slim:%s", uuid), "true", *conf.Cache.SkinCacheDuration); err != nil {
				return nil, false, err
			}
		} else {
			if err = r.Delete(fmt.Sprintf("slim:%s", uuid)); err != nil {
				return nil, false, err
			}
		}
	}

	if conf.Environment == "development" {
		log.Printf("Fetched player skin from Mojang (uuid=%s, slim=%v)\n", uuid, isSlim)
	}

	return skinImage, isSlim, nil
}

// EncodePNG encodes the image into PNG format and returns the data as a byte array.
func EncodePNG(img image.Image) ([]byte, error) {
	buf := &bytes.Buffer{}

	if err := png.Encode(buf, img); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// ParseQueryParams parses the query parameters from the request and returns a QueryParams struct, using default values from the provided configuration.
func ParseQueryParams(ctx *fiber.Ctx, route RouteConfig) *QueryParams {
	args := ctx.Context().QueryArgs()

	response := &QueryParams{
		Scale:    route.DefaultScale,
		Download: route.DefaultDownload,
		Overlay:  route.DefaultOverlay,
	}

	if args.Has("scale") {
		if scale, err := args.GetUint("scale"); err == nil {
			response.Scale = Clamp(scale, route.MinScale, route.MaxScale)
		}
	}

	if args.Has("overlay") {
		response.Overlay = args.GetBool("overlay")
	}

	if args.Has("download") {
		response.Download = args.GetBool("download")
	}

	return response
}

// GetInstanceID returns the INSTANCE_ID environment variable parsed as an unsigned 16-bit integer.
func GetInstanceID() (uint16, error) {
	if instanceID := os.Getenv("INSTANCE_ID"); len(instanceID) > 0 {
		value, err := strconv.ParseUint(instanceID, 10, 16)

		if err != nil {
			log.Fatal(err)
		}

		return uint16(value), nil
	}

	return 0, nil
}
