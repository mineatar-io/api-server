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
		AuthKey:     "",
		Redis: RedisConfig{
			Host:     "127.0.0.1",
			Port:     6379,
			User:     "",
			Password: "",
			Database: 0,
		},
		Routes: RoutesConfig{
			Face: RouteConfig{
				DefaultOverlay:  true,
				DefaultDownload: false,
				DefaultFallback: true,
				DefaultScale:    4,
				MinScale:        1,
				MaxScale:        64,
			},
			Head: RouteConfig{
				DefaultOverlay:  true,
				DefaultDownload: false,
				DefaultFallback: true,
				DefaultScale:    4,
				MinScale:        1,
				MaxScale:        64,
			},
			FullBody: RouteConfig{
				DefaultOverlay:  true,
				DefaultDownload: false,
				DefaultFallback: true,
				DefaultScale:    4,
				MinScale:        1,
				MaxScale:        64,
			},
			FrontBody: RouteConfig{
				DefaultOverlay:  true,
				DefaultDownload: false,
				DefaultFallback: true,
				DefaultScale:    4,
				MinScale:        1,
				MaxScale:        64,
			},
			BackBody: RouteConfig{
				DefaultOverlay:  true,
				DefaultDownload: false,
				DefaultFallback: true,
				DefaultScale:    4,
				MinScale:        1,
				MaxScale:        64,
			},
			LeftBody: RouteConfig{
				DefaultOverlay:  true,
				DefaultDownload: false,
				DefaultFallback: true,
				DefaultScale:    4,
				MinScale:        1,
				MaxScale:        64,
			},
			RightBody: RouteConfig{
				DefaultOverlay:  true,
				DefaultDownload: false,
				DefaultFallback: true,
				DefaultScale:    4,
				MinScale:        1,
				MaxScale:        64,
			},
			RawSkin: RouteConfig{
				DefaultDownload: false,
				DefaultFallback: true,
			},
		},
		Cache: CacheConfig{
			UUIDCacheDuration:   time.Hour * 168,
			SkinCacheDuration:   time.Hour * 12,
			RenderCacheDuration: time.Hour * 12,
		},
	}
)

// Config is the root configuration object for the application.
type Config struct {
	Environment string       `yaml:"environment"`
	Host        string       `yaml:"host"`
	Port        uint16       `yaml:"port"`
	AuthKey     string       `yaml:"auth_key"`
	Redis       RedisConfig  `yaml:"redis"`
	Routes      RoutesConfig `yaml:"routes"`
	Cache       CacheConfig  `yaml:"cache"`
}

// RoutesConfig is the configuration data of all API routes.
type RoutesConfig struct {
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
	DefaultFallback bool `yaml:"default_fallback"`
	DefaultDownload bool `yaml:"default_download"`
	MinScale        int  `yaml:"min_scale"`
	MaxScale        int  `yaml:"max_scale"`
}

// RedisConfig is the configuration data used to connect to Redis.
type RedisConfig struct {
	Host     string `yaml:"host"`
	Port     uint16 `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database int    `yaml:"database"`
}

// CacheConfig is the configuration data used to set TTL values for Redis keys.
type CacheConfig struct {
	UUIDCacheDuration   time.Duration `yaml:"uuid_cache_duration"`
	SkinCacheDuration   time.Duration `yaml:"skin_cache_duration"`
	RenderCacheDuration time.Duration `yaml:"render_cache_duration"`
}

// ReadFile reads the configuration from the file and parses it as YAML.
func (c *Config) ReadFile(file string) error {
	data, err := os.ReadFile(file)

	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, c)
}
