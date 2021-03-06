package conf

import (
	"io/ioutil"
	"time"

	"gopkg.in/yaml.v3"
)

type RouteConfig struct {
	DefaultScale    int  `yaml:"default_scale"`
	DefaultOverlay  bool `yaml:"default_overlay"`
	DefaultFallback bool `yaml:"default_fallback"`
	DefaultDownload bool `yaml:"default_download"`
	MinScale        int  `yaml:"min_scale"`
	MaxScale        int  `yaml:"max_scale"`
}

type Configuration struct {
	Redis struct {
		URI      string `yaml:"uri"`
		Database int    `yaml:"database"`
	} `yaml:"redis"`
	Routes struct {
		Face      RouteConfig `yaml:"face"`
		Head      RouteConfig `yaml:"head"`
		FullBody  RouteConfig `yaml:"full_body"`
		FrontBody RouteConfig `yaml:"front_body"`
		BackBody  RouteConfig `yaml:"back_body"`
		LeftBody  RouteConfig `yaml:"left_body"`
		RightBody RouteConfig `yaml:"right_body"`
		RawSkin   RouteConfig `yaml:"raw_skin"`
	} `yaml:"routes"`
	Cache struct {
		UUIDCacheDuration   time.Duration `yaml:"uuid_cache_duration"`
		SkinCacheDuration   time.Duration `yaml:"skin_cache_duration"`
		RenderCacheDuration time.Duration `yaml:"render_cache_duration"`
	} `yaml:"cache"`
}

func (c *Configuration) ReadFile(file string) error {
	data, err := ioutil.ReadFile(file)

	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, c)
}
