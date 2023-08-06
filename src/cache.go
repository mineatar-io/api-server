package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"image"
	"log"
)

type ResultCacheKey struct {
	UUID    string `json:"uuid"`
	Type    string `json:"type"`
	Scale   int    `json:"scale"`
	Overlay bool   `json:"overlay"`
}

// GetCacheKey returns the key used in the cache based on the rendering options, calculated as an SHA-256 hash.
func GetResultCacheKey(uuid, renderType string, opts *QueryParams) (string, error) {
	rawKey := ResultCacheKey{
		UUID:    uuid,
		Type:    renderType,
		Scale:   opts.Scale,
		Overlay: opts.Overlay,
	}

	rawKeyData, err := json.Marshal(rawKey)

	if err != nil {
		return "", err
	}

	hash := sha256.Sum256(rawKeyData)

	return fmt.Sprintf("result:%s", hex.EncodeToString(hash[:])), nil
}

// GetCachedRenderResult returns the render result from Redis cache, or nil if it does not exist or cache is disabled.
func GetCachedRenderResult(renderType, uuid string, opts *QueryParams) ([]byte, error) {
	if config.Cache.RenderCacheDuration == nil {
		return nil, nil
	}

	key, err := GetResultCacheKey(uuid, renderType, opts)

	if err != nil {
		return nil, err
	}

	return r.GetBytes(key)
}

// SetCachedRenderResult puts the render result into cache, or does nothing is cache is disabled.
func SetCachedRenderResult(renderType, uuid string, opts *QueryParams, data []byte) error {
	if config.Cache.RenderCacheDuration == nil {
		return nil
	}

	key, err := GetResultCacheKey(uuid, renderType, opts)

	if err != nil {
		return err
	}

	return r.Set(key, data, *config.Cache.RenderCacheDuration)
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

		if config.Environment == "development" {
			log.Printf("Retrieved player skin from cache (uuid=%s, slim=%v)\n", uuid, slim)
		}

		return cache, slim, nil
	}

	return nil, false, nil
}
