package routes

import (
	"github.com/mineatar-io/api-server/src/conf"
	"github.com/mineatar-io/api-server/src/redis"
)

var (
	r      *redis.Redis
	config *conf.Configuration
)

func Init(red *redis.Redis, c *conf.Configuration) {
	r = red
	config = c
}
