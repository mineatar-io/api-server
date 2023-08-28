package main

import (
	"fmt"
	"image"
	"net/url"
	"strconv"
)

type ResultCacheKey struct {
	UUID    string `json:"uuid"`
	Type    string `json:"type"`
	Scale   int    `json:"scale"`
	Overlay bool   `json:"overlay"`
}

// GetCacheKey returns the key used in the cache based on the rendering options, calculated as an SHA-256 hash.
func GetResultCacheKey(uuid, renderType string, opts *QueryParams) string {
	values := &url.Values{}
	values.Set("uuid", uuid)
	values.Set("type", renderType)
	values.Set("scale", strconv.FormatInt(int64(opts.Scale), 10))
	values.Set("overlay", strconv.FormatBool(opts.Overlay))
	values.Set("format", opts.Format)

	return SHA256(values.Encode())
}

// GetCachedRenderResult returns the render result from Redis cache, or nil if it does not exist or cache is disabled.
func GetCachedRenderResult(renderType, uuid string, opts *QueryParams) ([]byte, error) {
	if config.Cache.RenderCacheDuration == nil {
		return nil, nil
	}

	return r.GetBytes(fmt.Sprintf("result:%s", GetResultCacheKey(uuid, renderType, opts)))
}

// SetCachedRenderResult puts the render result into cache, or does nothing is cache is disabled.
func SetCachedRenderResult(renderType, uuid string, opts *QueryParams, data []byte) error {
	if config.Cache.RenderCacheDuration == nil {
		return nil
	}

	return r.Set(fmt.Sprintf("result:%s", GetResultCacheKey(uuid, renderType, opts)), data, *config.Cache.RenderCacheDuration)
}

// GetCachedSkin returns the raw skin of a player by UUID from the cache, also returning if the player has a slim player model.
func GetCachedSkin(uuid string) (*image.NRGBA, bool, error) {
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

	return nil, false, nil
}

func SetCachedSkin(uuid string, value []byte, isSlim bool) error {
	if err := r.Set(fmt.Sprintf("skin:%s", uuid), value, *config.Cache.SkinCacheDuration); err != nil {
		return err
	}

	if isSlim {
		if err := r.Set(fmt.Sprintf("slim:%s", uuid), "true", *config.Cache.SkinCacheDuration); err != nil {
			return err
		}
	} else {
		if err := r.Delete(fmt.Sprintf("slim:%s", uuid)); err != nil {
			return err
		}
	}

	return nil
}
