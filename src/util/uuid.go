package util

import (
	"fmt"
	"log"
	"strings"
	"time"
)

func GetUUID(value string) (string, error) {
	value = strings.ToLower(strings.ReplaceAll(value, "-", ""))

	if len(value) == 32 {
		return value, nil
	}

	cache, ok, err := r.GetString(fmt.Sprintf("uuid:%s", value))

	if err != nil {
		return "", err
	}

	if ok {
		if Debug {
			log.Printf("[DEBUG] Retrieved UUID for player %s (%s) from cache\n", value, cache)
		}

		return cache, nil
	}

	profile, err := GetProfileByUsername(value)

	if err != nil {
		return "", err
	}

	if profile == nil {
		return "", nil
	}

	if err = r.Set(fmt.Sprintf("uuid:%s", value), profile.ID, time.Duration(config.Cache.UUIDCacheDuration)*time.Second); err != nil {
		return "", err
	}

	if Debug {
		log.Printf("[DEBUG] Fetched UUID for player %s (%s) from Mojang\n", profile.Name, profile.ID)
	}

	return profile.ID, nil
}
