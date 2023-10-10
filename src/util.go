package main

import (
	"bytes"
	"crypto/sha256"
	_ "embed"
	"encoding/hex"
	"errors"
	"fmt"
	"image"
	"image/draw"
	"image/gif"
	"image/jpeg"
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
	faviconData    []byte
	AllowedFormats []string = []string{
		"png",
		"jpg",
		"jpeg",
		"gif",
	}
)

// QueryParams is used by most all API routes as options for how the image should be rendered, or how errors should be handled.
type QueryParams struct {
	Scale    int
	Download bool
	Overlay  bool
	Format   string
}

// PointerOf returns the value of the first argument as a pointer.
func PointerOf[T any](v T) *T {
	return &v
}

// Contains returns true if the array contains the value.
func Contains[T comparable](arr []T, value T) bool {
	for _, v := range arr {
		if v != value {
			continue
		}

		return true
	}

	return false
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
	if config.Cache.EnableLocks {
		mutex := r.NewMutex(fmt.Sprintf("skin-lock:%s", uuid))
		mutex.Lock()

		defer mutex.Unlock()
	}

	// Get skin from cache, and return if it exists
	if config.Cache.SkinCacheDuration != nil {
		rawSkin, slim, err := GetCachedSkin(uuid)

		if err != nil {
			return nil, false, err
		}

		if rawSkin != nil {
			return rawSkin, slim, nil
		}
	}

	var (
		err              error             = nil
		skinImage        *image.NRGBA      = nil
		rawSkin          []byte            = nil
		isSlim           bool              = skin.IsSlimFromUUID(uuid)
		profile          *MinecraftProfile = nil
		rawTextures      string            = ""
		texturesProperty *DecodedTextures  = nil
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
		if texturesProperty, err = DecodeTexturesValue(rawTextures); err != nil {
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
			if !errors.Is(err, image.ErrFormat) {
				return nil, false, err
			}

			skinImage = skin.GetDefaultSkin(isSlim)
		}

		if rawSkin, err = EncodePNG(skinImage); err != nil {
			return nil, false, err
		}
	}

	// Put the skin into cache so it can be used for future requests
	if config.Cache.SkinCacheDuration != nil {
		if err = SetCachedSkin(uuid, rawSkin, isSlim); err != nil {
			return nil, false, err
		}
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

// EncodeImage encodes the image into the format specified by the query parameters.
func EncodeImage(img image.Image, opts *QueryParams) ([]byte, error) {
	buf := &bytes.Buffer{}

	switch opts.Format {
	case "png":
		{
			if err := png.Encode(buf, img); err != nil {
				return nil, err
			}

			break
		}
	case "jpg", "jpeg":
		{
			if err := jpeg.Encode(buf, img, nil); err != nil {
				return nil, err
			}

			break
		}
	case "gif":
		{
			if err := gif.Encode(buf, img, nil); err != nil {
				return nil, err
			}

			break
		}
	default:
		return nil, fmt.Errorf("invalid format: %s", opts.Format)
	}

	return buf.Bytes(), nil
}

// ParseQueryParams parses the query parameters from the request and returns a QueryParams struct, using default values from the provided configuration.
func ParseQueryParams(ctx *fiber.Ctx, route RouteConfig) *QueryParams {
	format := ctx.Query("format", route.DefaultFormat)

	if !Contains(AllowedFormats, format) {
		ctx.Status(http.StatusBadRequest).SendString("Invalid 'format' query parameter")

		return nil
	}

	return &QueryParams{
		Scale:    Clamp(ctx.QueryInt("scale", route.DefaultScale), route.MinScale, route.MaxScale),
		Download: ctx.QueryBool("download", route.DefaultDownload),
		Overlay:  ctx.QueryBool("overlay", route.DefaultOverlay),
		Format:   ctx.Query("format", route.DefaultFormat),
	}
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

// SHA256 computes the SHA-256 hash of the input string.
func SHA256(value string) string {
	hash := sha256.Sum256([]byte(value))

	return hex.EncodeToString(hash[:])
}
