package main

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

var (
	// DefaultConfig is the default configuration values used by the application.
	DefaultConfig *Config = &Config{
		Environment: "development",
		Host:        "127.0.0.1",
		Port:        3001,
		Redis:       "redis://127.0.0.1:6379/0",
		Routes: Routes{
			Face: RouteConfig{
				DefaultOverlay:  true,
				DefaultDownload: false,
				DefaultScale:    4,
				MinScale:        1,
				MaxScale:        64,
			},
			Head: RouteConfig{
				DefaultOverlay:  true,
				DefaultDownload: false,
				DefaultScale:    4,
				MinScale:        1,
				MaxScale:        64,
			},
			FullBody: RouteConfig{
				DefaultOverlay:  true,
				DefaultDownload: false,
				DefaultScale:    4,
				MinScale:        1,
				MaxScale:        64,
			},
			FrontBody: RouteConfig{
				DefaultOverlay:  true,
				DefaultDownload: false,
				DefaultScale:    4,
				MinScale:        1,
				MaxScale:        64,
			},
			BackBody: RouteConfig{
				DefaultOverlay:  true,
				DefaultDownload: false,
				DefaultScale:    4,
				MinScale:        1,
				MaxScale:        64,
			},
			LeftBody: RouteConfig{
				DefaultOverlay:  true,
				DefaultDownload: false,
				DefaultScale:    4,
				MinScale:        1,
				MaxScale:        64,
			},
			RightBody: RouteConfig{
				DefaultOverlay:  true,
				DefaultDownload: false,
				DefaultScale:    4,
				MinScale:        1,
				MaxScale:        64,
			},
			RawSkin: RouteConfig{
				DefaultDownload: false,
			},
		},
		Cache: CacheConfig{
			SkinCacheDuration:   PointerOf(time.Hour * 12),
			RenderCacheDuration: PointerOf(time.Hour * 12),
			EnableLocks:         true,
		},
	}
)

// Config is the root configuration object for the application.
type Config struct {
	Environment string      `yaml:"environment"`
	Host        string      `yaml:"host"`
	Port        uint16      `yaml:"port"`
	Redis       string      `yaml:"redis"`
	Routes      Routes      `yaml:"routes"`
	Cache       CacheConfig `yaml:"cache"`
}

// Routes is the configuration data of all API routes.
type Routes struct {
	Face      RouteConfig `yaml:"face"`
	Head      RouteConfig `yaml:"head"`
	FullBody  RouteConfig `yaml:"full_body"`
	FrontBody RouteConfig `yaml:"front_body"`
	BackBody  RouteConfig `yaml:"back_body"`
	LeftBody  RouteConfig `yaml:"left_body"`
	RightBody RouteConfig `yaml:"right_body"`
	RawSkin   RouteConfig `yaml:"raw_skin"`
}

// RouteConfig is the configuration data used by a single API route.
type RouteConfig struct {
	DefaultScale    int  `yaml:"default_scale"`
	DefaultOverlay  bool `yaml:"default_overlay"`
	DefaultDownload bool `yaml:"default_download"`
	MinScale        int  `yaml:"min_scale"`
	MaxScale        int  `yaml:"max_scale"`
}

// CacheConfig is the configuration data used to set TTL values for Redis keys.
type CacheConfig struct {
	SkinCacheDuration   *time.Duration `yaml:"skin_cache_duration"`
	RenderCacheDuration *time.Duration `yaml:"render_cache_duration"`
	EnableLocks         bool           `yaml:"enable_locks"`
}

// ReadFile reads the configuration from the file and parses it as YAML.
func (c *Config) ReadFile(file string) error {
	data, err := os.ReadFile(file)

	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, c)
}
