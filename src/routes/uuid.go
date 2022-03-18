package routes

import (
	"log"
	"net/http"

	"github.com/mineatar-io/api-server/src/util"
	"github.com/valyala/fasthttp"
)

func UUIDHandler(ctx *fasthttp.RequestCtx) {
	user := ctx.UserValue("user").(string)

	uuid, ok, err := util.LookupUUID(user)

	if err != nil {
		log.Println(err)

		ctx.SetStatusCode(http.StatusInternalServerError)
		ctx.SetBodyString(http.StatusText(http.StatusInternalServerError))

		return
	}

	if !ok {
		ctx.SetStatusCode(404)
		ctx.SetBodyString(http.StatusText(http.StatusNotFound))

		return
	}

	ctx.SetBodyString(uuid)
}
