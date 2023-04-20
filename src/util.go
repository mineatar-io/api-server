package main

import (
	"bytes"
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

// QueryParams is used by most all API routes as options for how the image should be rendered, or how errors should be handled.
type QueryParams struct {
	Scale    int  `query:"scale"`
	Download bool `query:"download"`
	Overlay  bool `query:"overlay"`
	Fallback bool `query:"fallback"`
}

// FormatUUID returns the UUID string without any dashes.
func FormatUUID(uuid string) string {
	return strings.ToLower(strings.ReplaceAll(uuid, "-", ""))
}

// LookupUUID returns the UUID of a player either by username or UUID, while attempting to use any cached values in the database.
func LookupUUID(value string) (string, bool, error) {
	value = FormatUUID(value)

	if len(value) == 32 {
		return value, true, nil
	}

	cacheKey := fmt.Sprintf("uuid:%s", value)

	cache, ok, err := r.GetString(cacheKey)

	if err != nil {
		return "", false, err
	}

	if ok {
		return cache, true, nil
	}

	profile, err := UsernameToUUID(value)

	if err != nil {
		return "", false, err
	}

	if profile == nil {
		return "", false, nil
	}

	if err = r.Set(cacheKey, profile.UUID, conf.Cache.UUIDCacheDuration); err != nil {
		return "", true, err
	}

	return profile.UUID, true, nil
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
	uuid = FormatUUID(uuid)

	if len(uuid) < 1 {
		return skin.GetDefaultSkin(false), false, nil
	}

	cache, ok, err := r.GetNRGBA(fmt.Sprintf("skin:%s", uuid))

	if err != nil {
		return nil, false, err
	}

	if ok {
		slim, err := r.Exists(fmt.Sprintf("slim:%s", uuid))

		if err != nil {
			return nil, false, err
		}

		return cache, slim, nil
	}

	textures, err := GetProfileTextures(uuid)

	if err != nil {
		return nil, false, err
	}

	if textures == nil {
		slim := skin.IsSlimFromUUID(uuid)

		return skin.GetDefaultSkin(slim), slim, nil
	}

	value := ""

	for _, property := range textures.Properties {
		if property.Name != "textures" {
			continue
		}

		value = property.Value
	}

	if len(value) < 1 {
		slim := skin.IsSlimFromUUID(uuid)

		return skin.GetDefaultSkin(slim), slim, nil
	}

	texturesResult, err := GetDecodedTexturesValue(value)

	if err != nil {
		return nil, false, err
	}

	if len(texturesResult.Textures.Skin.URL) < 1 {
		slim := skin.IsSlimFromUUID(uuid)

		return skin.GetDefaultSkin(slim), slim, nil
	}

	slim := texturesResult.Textures.Skin.Metadata.Model == "slim"

	skin, err := FetchImage(texturesResult.Textures.Skin.URL)

	if err != nil {
		return nil, false, err
	}

	encodedSkin, err := EncodePNG(skin)

	if err != nil {
		return nil, false, err
	}

	if err = r.Set(fmt.Sprintf("skin:%s", uuid), encodedSkin, conf.Cache.SkinCacheDuration); err != nil {
		return nil, false, err
	}

	if slim {
		if err = r.Set(fmt.Sprintf("slim:%s", uuid), "true", conf.Cache.SkinCacheDuration); err != nil {
			return nil, false, err
		}
	} else {
		if err = r.Delete(fmt.Sprintf("slim:%s", uuid)); err != nil {
			return nil, false, err
		}
	}

	return skin, slim, nil
}

// Clamp clamps the input value between the minimum and maximum values.
// This method is preferred over `math.Min()` and `math.Max()` to prevent any type coercion between floats and integers.
func Clamp(value, min, max int) int {
	if value > max {
		return max
	}

	if value < min {
		return min
	}

	return value
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
		Fallback: route.DefaultFallback,
	}

	if args.Has("scale") {
		if scale, err := args.GetUint("scale"); err == nil {
			response.Scale = Clamp(scale, route.MinScale, route.MaxScale)
		}
	}

	if args.Has("overlay") {
		response.Overlay = args.GetBool("overlay")
	}

	if args.Has("fallback") {
		response.Fallback = args.GetBool("fallback")
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

// ParseUserParam returns the user name from the route param, allowing values such as "PassTheMayo.png" to be returned as "PassTheMayo".
func ParseUserParam(ctx *fiber.Ctx) string {
	return strings.Split(ctx.Params("user"), ".")[0]
}
