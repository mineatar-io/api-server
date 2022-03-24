package util

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/mineatar-io/api-server/src/conf"
	"github.com/mineatar-io/skin-render"
	"github.com/mineatar-io/yggdrasil"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/valyala/fasthttp"
)

var (
	Debug                     bool = os.Getenv("DEBUG") == "true"
	yggdrasilUUIDLookupMetric      = promauto.NewCounter(prometheus.CounterOpts{
		Name: "yggdrasil_uuid_lookup_count",
		Help: "The amount of Yggdrasil UUID lookup requests",
	})
	yggdrasilTextureLookupMetric = promauto.NewCounter(prometheus.CounterOpts{
		Name: "yggdrasil_texture_lookup_count",
		Help: "The amount of Yggdrasil texture lookup requests",
	})
)

type QueryParams struct {
	Scale    int
	Download bool
	Overlay  bool
	Fallback bool
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
		if Debug {
			log.Printf("Retrieved UUID from cache for '%s' (%s)\n", value, cache)
		}

		return cache, true, nil
	}

	profile, err := yggdrasil.UsernameToUUID(value)

	if err != nil {
		return "", false, err
	}

	yggdrasilUUIDLookupMetric.Inc()

	if profile == nil {
		if Debug {
			log.Printf("Fetched UUID from Mojang for '%s', did not exist\n", value)
		}

		return "", false, nil
	}

	if err = r.Set(cacheKey, profile.UUID, config.Cache.UUIDCacheDuration); err != nil {
		return "", true, err
	}

	if Debug {
		log.Printf("Fetched UUID from Mojang for '%s' (%s)\n", value, profile.UUID)
	}

	return profile.UUID, true, nil
}

func FetchImage(url string) (*image.NRGBA, error) {
	resp, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	if Debug {
		log.Printf("Fetched image from URL: %s\n", url)
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

		if Debug {
			log.Printf("Retrieved skin for '%s' (slim: %t) from cache\n", uuid, slim)
		}

		return cache, slim, nil
	}

	textures, err := yggdrasil.GetProfileTextures(uuid)

	if err != nil {
		return nil, false, err
	}

	yggdrasilTextureLookupMetric.Inc()

	if textures == nil {
		if Debug {
			log.Printf("Fetched textures for '%s' from Mojang, none exists, using default skin\n", uuid)
		}

		slim := skin.IsSlimFromUUID(uuid)

		return skin.GetDefaultSkin(slim), slim, nil
	}

	if Debug {
		log.Printf("Fetched textures for '%s' from Mojang\n", uuid)
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

	if err = r.Set(fmt.Sprintf("skin:%s", uuid), encodedSkin, time.Duration(config.Cache.SkinCacheDuration)*time.Second); err != nil {
		return nil, false, err
	}

	if slim {
		if err = r.Set(fmt.Sprintf("slim:%s", uuid), "true", time.Duration(config.Cache.SkinCacheDuration)*time.Second); err != nil {
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

func ParseQueryParams(ctx *fasthttp.RequestCtx, route conf.RouteConfig) *QueryParams {
	args := ctx.QueryArgs()

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
