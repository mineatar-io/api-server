package conf

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Configuration struct {
	Redis struct {
		URI      string `yaml:"uri"`
		Database int    `yaml:"database"`
	} `yaml:"redis"`
	Routes struct {
		Face struct {
			DefaultScale int `yaml:"default_scale"`
			MinScale     int `yaml:"min_scale"`
			MaxScale     int `yaml:"max_scale"`
		} `yaml:"face"`
		Head struct {
			DefaultScale int `yaml:"default_scale"`
			MinScale     int `yaml:"min_scale"`
			MaxScale     int `yaml:"max_scale"`
		} `yaml:"head"`
		FullBody struct {
			DefaultScale int `yaml:"default_scale"`
			MinScale     int `yaml:"min_scale"`
			MaxScale     int `yaml:"max_scale"`
		} `yaml:"full_body"`
		FrontBody struct {
			DefaultScale int `yaml:"default_scale"`
			MinScale     int `yaml:"min_scale"`
			MaxScale     int `yaml:"max_scale"`
		} `yaml:"front_body"`
		BackBody struct {
			DefaultScale int `yaml:"default_scale"`
			MinScale     int `yaml:"min_scale"`
			MaxScale     int `yaml:"max_scale"`
		} `yaml:"back_body"`
		LeftBody struct {
			DefaultScale int `yaml:"default_scale"`
			MinScale     int `yaml:"min_scale"`
			MaxScale     int `yaml:"max_scale"`
		} `yaml:"left_body"`
		RightBody struct {
			DefaultScale int `yaml:"default_scale"`
			MinScale     int `yaml:"min_scale"`
			MaxScale     int `yaml:"max_scale"`
		} `yaml:"right_body"`
	} `yaml:"routes"`
	Cache struct {
		UUIDCacheDuration   int64 `yaml:"uuid_cache_duration"`
		SkinCacheDuration   int64 `yaml:"skin_cache_duration"`
		RenderCacheDuration int64 `yaml:"render_cache_duration"`
	} `yaml:"cache"`
}

func (c *Configuration) ReadFile(file string) error {
	data, err := ioutil.ReadFile(file)

	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, c)
}
