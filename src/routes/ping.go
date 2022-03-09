package routes

import (
	"github.com/valyala/fasthttp"
)

func PingHandler(ctx *fasthttp.RequestCtx) {
	ctx.SetBodyString("Pong!")
}
