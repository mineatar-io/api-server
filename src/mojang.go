package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// MinecraftProfile is metadata about a Minecraft player returned from the Mojang API.
type MinecraftProfile struct {
	UUID       string `json:"id"`
	Username   string `json:"name"`
	Legacy     bool   `json:"legacy"`
	Properties []struct {
		Name      string `json:"name"`
		Value     string `json:"value"`
		Signature string `json:"signature,omitempty"`
	} `json:"properties"`
}

// DecodedTextures is the decoded object of the base64-encoded values property in a MinecraftProfile properties value.
type DecodedTextures struct {
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

// GetMinecraftProfile returns the textures of a Minecraft player from Mojang.
func GetMinecraftProfile(uuid string) (*MinecraftProfile, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://sessionserver.mojang.com/session/minecraft/profile/%s", uuid), nil)

	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "mineatar.io")

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

	var response MinecraftProfile

	if err = json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// DecodeTexturesValue decodes the value from a MinecraftProfile texture property.
func DecodeTexturesValue(value string) (*DecodedTextures, error) {
	rawResult, err := base64.StdEncoding.DecodeString(value)

	if err != nil {
		return nil, err
	}

	var result DecodedTextures

	if err = json.Unmarshal(rawResult, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
