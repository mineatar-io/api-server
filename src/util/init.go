package util

import (
	"os"

	"github.com/mineatar-io/api-server/src/conf"
	"github.com/mineatar-io/api-server/src/redis"
)

var (
	Debug  = os.Getenv("DEBUG") == "true"
	config *conf.Configuration
	r      *redis.Redis
)

func Init(red *redis.Redis, c *conf.Configuration) {
	r = red
	config = c
}
