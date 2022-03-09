package util

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"image/draw"
	"io/ioutil"
	"log"
	"main/src/redis"
	"net/http"
	"time"
)

type BasicProfile struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

type Profile struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Properties []struct {
		Name      string `json:"name"`
		Value     string `json:"value"`
		Signature string `json:"signature"`
	} `yaml:"properties"`
}

type Textures struct {
	Textures map[string]struct {
		URL      string `json:"url"`
		Metadata struct {
			Model string `json:"model"`
		} `json:"metadata"`
	} `json:"textures"`
}

func GetProfileByUsername(username string) (*BasicProfile, error) {
	resp, err := http.Get(fmt.Sprintf("https://api.mojang.com/users/profiles/minecraft/%s", username))

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		if resp.StatusCode == 204 {
			return nil, nil
		}

		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	result := &BasicProfile{}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(body, result); err != nil {
		return nil, err
	}

	return result, nil
}

func GetPlayerProfile(uuid string) (*Profile, error) {
	resp, err := http.Get(fmt.Sprintf("https://sessionserver.mojang.com/session/minecraft/profile/%s", uuid))

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		if resp.StatusCode == 404 {
			return nil, nil
		}

		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	result := &Profile{}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(body, result); err != nil {
		return nil, err
	}

	return result, nil
}

func GetPlayerSkin(r *redis.Redis, uuid string) (*image.NRGBA, bool, error) {
	if len(uuid) < 1 {
		return GetDefaultSkin(false), false, nil
	}

	cache, ok, err := r.GetBytes(fmt.Sprintf("skin:%s", uuid))

	if err != nil {
		return nil, false, err
	}

	if ok {
		slim, err := r.Exists(fmt.Sprintf("slim:%s", uuid))

		if err != nil {
			return nil, false, err
		}

		img, format, err := image.Decode(bytes.NewReader(cache))

		if err != nil {
			return nil, false, err
		}

		if Debug {
			log.Printf("[DEBUG] Retrieved skin for %s from cache\n", uuid)
		}

		if format != "NRGBA" {
			newImage := image.NewNRGBA(img.Bounds())

			draw.Draw(newImage, img.Bounds(), img, image.Pt(0, 0), draw.Src)

			return newImage, slim, nil
		}

		return img.(*image.NRGBA), slim, nil
	}

	profile, err := GetPlayerProfile(uuid)

	if err != nil {
		return nil, false, err
	}

	if profile == nil {
		return nil, false, nil
	}

	var value string

	for _, property := range profile.Properties {
		if property.Name != "textures" {
			continue
		}

		value = property.Value
	}

	if len(value) < 1 {
		return nil, false, nil
	}

	propertyValue, err := base64.StdEncoding.DecodeString(value)

	if err != nil {
		return nil, false, err
	}

	textures := &Textures{}

	if err = json.Unmarshal(propertyValue, textures); err != nil {
		return nil, false, err
	}

	skinProperty, ok := textures.Textures["SKIN"]

	if !ok {
		return nil, false, nil
	}

	slim := skinProperty.Metadata.Model == "slim"

	resp, err := http.Get(skinProperty.URL)

	if err != nil {
		return nil, false, err
	}

	if resp.StatusCode != 200 {
		if resp.StatusCode == 404 {
			return nil, slim, nil
		}

		return nil, false, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, false, err
	}

	img, format, err := image.Decode(bytes.NewReader(body))

	if err != nil {
		return nil, false, err
	}

	if slim {
		if err = r.Set(fmt.Sprintf("slim:%s", uuid), "true", time.Hour*24); err != nil {
			return nil, false, err
		}
	} else {
		if err = r.Delete(fmt.Sprintf("slim:%s", uuid)); err != nil {
			return nil, false, err
		}
	}

	if err = r.Set(fmt.Sprintf("skin:%s", uuid), body, time.Hour*24); err != nil {
		return nil, false, err
	}

	if Debug {
		log.Printf("[DEBUG] Fetched skin for %s from Mojang\n", uuid)
	}

	if format != "NRGBA" {
		newImage := image.NewNRGBA(img.Bounds())

		draw.Draw(newImage, img.Bounds(), img, image.Pt(0, 0), draw.Src)

		return newImage, slim, nil
	}

	return img.(*image.NRGBA), slim, nil
}
