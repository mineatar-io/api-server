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
	"github.com/mineatar-io/yggdrasil"
)

type QueryParams struct {
	Scale    int  `query:"scale"`
	Download bool `query:"download"`
	Overlay  bool `query:"overlay"`
	Fallback bool `query:"fallback"`
}

func FormatUUID(uuid string) string {
	return strings.ToLower(strings.ReplaceAll(uuid, "-", ""))
}

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

	profile, err := yggdrasil.UsernameToUUID(value)

	if err != nil {
		return "", false, err
	}

	if profile == nil {
		return "", false, nil
	}

	if err = r.Set(cacheKey, profile.UUID, config.Cache.UUIDCacheDuration); err != nil {
		return "", true, err
	}

	return profile.UUID, true, nil
}

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

	textures, err := yggdrasil.GetProfileTextures(uuid)

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

	texturesResult, err := yggdrasil.GetDecodedTexturesValue(value)

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

	if err = r.Set(fmt.Sprintf("skin:%s", uuid), encodedSkin, config.Cache.SkinCacheDuration); err != nil {
		return nil, false, err
	}

	if slim {
		if err = r.Set(fmt.Sprintf("slim:%s", uuid), "true", config.Cache.SkinCacheDuration); err != nil {
			return nil, false, err
		}
	} else {
		if err = r.Delete(fmt.Sprintf("slim:%s", uuid)); err != nil {
			return nil, false, err
		}
	}

	return skin, slim, nil
}

// This is used instead of `math.Min/Max` because of the
// unnecessary coercion from/to float64.
func Clamp(value, min, max int) int {
	if value > max {
		return max
	}

	if value < min {
		return min
	}

	return value
}

func EncodePNG(img image.Image) ([]byte, error) {
	buf := &bytes.Buffer{}

	if err := png.Encode(buf, img); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

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

func ParseUserParam(ctx *fiber.Ctx) string {
	return strings.Split(ctx.Params("user"), ".")[0]
}
