package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// MinecraftProfile is Minecraft profile information returned from the Mojang API.
type MinecraftProfile struct {
	Username string `json:"name"`
	UUID     string `json:"id"`
}

// MinecraftProfileTextures is texture information about a Minecraft profile returned from the Mojang API.
type MinecraftProfileTextures struct {
	UUID       string `json:"id"`
	Username   string `json:"name"`
	Legacy     bool   `json:"legacy"`
	Properties []struct {
		Name      string `json:"name"`
		Value     string `json:"value"`
		Signature string `json:"signature,omitempty"`
	} `json:"properties"`
}

// MinecraftDecodedTextures is the decoded object of the Base64-encoded values property in a MinecraftProfileTextures texture value.
type MinecraftDecodedTextures struct {
	Timestamp         int64  `json:"timestamp"`
	UUID              string `json:"uuid"`
	Username          string `json:"username"`
	SignatureRequired bool   `json:"signatureRequired"`
	Textures          struct {
		Skin struct {
			URL      string `json:"url"`
			Metadata struct {
				Model string `json:"model"`
			} `json:"metadata,omitempty"`
		} `json:"SKIN,omitempty"`
		Cape struct {
			URL string `json:"url"`
		} `json:"CAPE,omitempty"`
	} `json:"textures"`
}

// UsernameToUUID converts a Minecraft username into a UUID using Mojang.
func UsernameToUUID(username string) (*MinecraftProfile, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.mojang.com/users/profiles/minecraft/%s", username), nil)

	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "mineatar.io Skin Render API")

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNoContent || resp.StatusCode == http.StatusNotFound {
			return nil, nil
		}

		return nil, fmt.Errorf("mojang: unexpected response: %s", resp.Status)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	response := &MinecraftProfile{}

	if err = json.Unmarshal(body, response); err != nil {
		return nil, err
	}

	return response, nil
}

// GetProfileTextures returns the textures of a Minecraft player from Mojang.
func GetProfileTextures(uuid string) (*MinecraftProfileTextures, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://sessionserver.mojang.com/session/minecraft/profile/%s", uuid), nil)

	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "mineatar.io Skin Render API")

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNoContent {
			return nil, nil
		}

		return nil, fmt.Errorf("mojang: unexpected response: %s", resp.Status)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	response := &MinecraftProfileTextures{}

	if err = json.Unmarshal(body, response); err != nil {
		return nil, err
	}

	return response, nil
}

// GetDecodedTexturesValue decodes the values from a MinecraftProfileTextures texture value.
func GetDecodedTexturesValue(value string) (*MinecraftDecodedTextures, error) {
	rawResult, err := base64.StdEncoding.DecodeString(value)

	if err != nil {
		return nil, err
	}

	result := MinecraftDecodedTextures{}

	if err = json.Unmarshal(rawResult, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
