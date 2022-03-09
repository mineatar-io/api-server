package routes

import "main/src/redis"

var (
	r *redis.Redis
)

const (
	MaxScale         float64 = 64.0
	MaxScaleFullBody float64 = 32.0
	MinScale         float64 = 1.0
)

func InitRoutes(red *redis.Redis) {
	r = red
}
