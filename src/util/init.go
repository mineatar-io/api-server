package util

import (
	"github.com/mineatar-io/api-server/src/conf"
	"github.com/mineatar-io/api-server/src/redis"
)

var (
	config *conf.Configuration
	r      *redis.Redis
)

func Init(red *redis.Redis, c *conf.Configuration) {
	r = red
	config = c
}
