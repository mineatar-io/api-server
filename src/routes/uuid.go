package routes

import (
	"log"
	"main/src/util"
	"net/http"

	"github.com/valyala/fasthttp"
)

func UUIDHandler(ctx *fasthttp.RequestCtx) {
	user := ctx.UserValue("user").(string)

	uuid, err := util.GetUUID(r, user)

	if err != nil {
		log.Println(err)

		ctx.SetStatusCode(http.StatusInternalServerError)
		ctx.SetBodyString(http.StatusText(http.StatusInternalServerError))

		return
	}

	ctx.SetBodyString(uuid)
}
